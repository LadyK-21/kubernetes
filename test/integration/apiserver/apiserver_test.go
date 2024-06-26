/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apiserver

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/features"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	utilfeaturetesting "k8s.io/apiserver/pkg/util/feature/testing"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/pager"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/api/testapi"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/master"
	"k8s.io/kubernetes/test/integration/framework"
)

func setup(t *testing.T, groupVersions ...schema.GroupVersion) (*httptest.Server, clientset.Interface, framework.CloseFunc) {
	masterConfig := framework.NewIntegrationTestMasterConfig()
	if len(groupVersions) > 0 {
		resourceConfig := master.DefaultAPIResourceConfigSource()
		resourceConfig.EnableVersions(groupVersions...)
		masterConfig.ExtraConfig.APIResourceConfigSource = resourceConfig
	}
	masterConfig.GenericConfig.OpenAPIConfig = framework.DefaultOpenAPIConfig()
	_, s, closeFn := framework.RunAMaster(masterConfig)

	clientSet, err := clientset.NewForConfig(&restclient.Config{Host: s.URL})
	if err != nil {
		t.Fatalf("Error in create clientset: %v", err)
	}
	return s, clientSet, closeFn
}

func verifyStatusCode(t *testing.T, verb, URL, body string, expectedStatusCode int) {
	// We dont use the typed Go client to send this request to be able to verify the response status code.
	bodyBytes := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(verb, URL, bodyBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v in sending req with verb: %s, URL: %s and body: %s", err, verb, URL, body)
	}
	transport := http.DefaultTransport
	klog.Infof("Sending request: %v", req)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v in req: %v", err, req)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected status %v, but got %v", expectedStatusCode, resp.StatusCode)
		t.Errorf("Body: %v", string(b))
	}
}

func path(resource, namespace, name string) string {
	return testapi.Extensions.ResourcePath(resource, namespace, name)
}

func newRS(namespace string) *apps.ReplicaSet {
	return &apps.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: "apiserver-test",
		},
		Spec: apps.ReplicaSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"name": "test"}},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"name": "test"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "fake-name",
							Image: "fakeimage",
						},
					},
				},
			},
		},
	}
}

var cascDel = `
{
  "kind": "DeleteOptions",
  "apiVersion": "` + testapi.Groups[api.GroupName].GroupVersion().String() + `",
  "orphanDependents": false
}
`

// Tests that the apiserver returns 202 status code as expected.
func Test202StatusCode(t *testing.T) {
	s, clientSet, closeFn := setup(t)
	defer closeFn()

	ns := framework.CreateTestingNamespace("status-code", s, t)
	defer framework.DeleteTestingNamespace(ns, s, t)

	rsClient := clientSet.AppsV1().ReplicaSets(ns.Name)

	// 1. Create the resource without any finalizer and then delete it without setting DeleteOptions.
	// Verify that server returns 200 in this case.
	rs, err := rsClient.Create(newRS(ns.Name))
	if err != nil {
		t.Fatalf("Failed to create rs: %v", err)
	}
	verifyStatusCode(t, "DELETE", s.URL+path("replicasets", ns.Name, rs.Name), "", 200)

	// 2. Create the resource with a finalizer so that the resource is not immediately deleted and then delete it without setting DeleteOptions.
	// Verify that the apiserver still returns 200 since DeleteOptions.OrphanDependents is not set.
	rs = newRS(ns.Name)
	rs.ObjectMeta.Finalizers = []string{"kube.io/dummy-finalizer"}
	rs, err = rsClient.Create(rs)
	if err != nil {
		t.Fatalf("Failed to create rs: %v", err)
	}
	verifyStatusCode(t, "DELETE", s.URL+path("replicasets", ns.Name, rs.Name), "", 200)

	// 3. Create the resource and then delete it with DeleteOptions.OrphanDependents=false.
	// Verify that the server still returns 200 since the resource is immediately deleted.
	rs = newRS(ns.Name)
	rs, err = rsClient.Create(rs)
	if err != nil {
		t.Fatalf("Failed to create rs: %v", err)
	}
	verifyStatusCode(t, "DELETE", s.URL+path("replicasets", ns.Name, rs.Name), cascDel, 200)

	// 4. Create the resource with a finalizer so that the resource is not immediately deleted and then delete it with DeleteOptions.OrphanDependents=false.
	// Verify that the server returns 202 in this case.
	rs = newRS(ns.Name)
	rs.ObjectMeta.Finalizers = []string{"kube.io/dummy-finalizer"}
	rs, err = rsClient.Create(rs)
	if err != nil {
		t.Fatalf("Failed to create rs: %v", err)
	}
	verifyStatusCode(t, "DELETE", s.URL+path("replicasets", ns.Name, rs.Name), cascDel, 202)
}

func TestAPIListChunking(t *testing.T) {
	defer utilfeaturetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.APIListChunking, true)()
	s, clientSet, closeFn := setup(t)
	defer closeFn()

	ns := framework.CreateTestingNamespace("list-paging", s, t)
	defer framework.DeleteTestingNamespace(ns, s, t)

	rsClient := clientSet.AppsV1().ReplicaSets(ns.Name)

	for i := 0; i < 4; i++ {
		rs := newRS(ns.Name)
		rs.Name = fmt.Sprintf("test-%d", i)
		if _, err := rsClient.Create(rs); err != nil {
			t.Fatal(err)
		}
	}

	calls := 0
	firstRV := ""
	p := &pager.ListPager{
		PageSize: 1,
		PageFn: pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
			calls++
			list, err := rsClient.List(opts)
			if err != nil {
				return nil, err
			}
			if calls == 1 {
				firstRV = list.ResourceVersion
			}
			if calls == 2 {
				rs := newRS(ns.Name)
				rs.Name = "test-5"
				if _, err := rsClient.Create(rs); err != nil {
					t.Fatal(err)
				}
			}
			return list, err
		}),
	}
	listObj, err := p.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 4 {
		t.Errorf("unexpected list invocations: %d", calls)
	}
	list := listObj.(metav1.ListInterface)
	if len(list.GetContinue()) != 0 {
		t.Errorf("unexpected continue: %s", list.GetContinue())
	}
	if list.GetResourceVersion() != firstRV {
		t.Errorf("unexpected resource version: %s instead of %s", list.GetResourceVersion(), firstRV)
	}
	var names []string
	if err := meta.EachListItem(listObj, func(obj runtime.Object) error {
		rs := obj.(*apps.ReplicaSet)
		names = append(names, rs.Name)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(names, []string{"test-0", "test-1", "test-2", "test-3"}) {
		t.Errorf("unexpected items: %#v", list)
	}
}

func makeSecret(name string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string][]byte{
			"key": []byte("value"),
		},
	}
}

func TestNameInFieldSelector(t *testing.T) {
	s, clientSet, closeFn := setup(t)
	defer closeFn()

	numNamespaces := 3
	namespaces := make([]*v1.Namespace, 0, numNamespaces)
	for i := 0; i < 3; i++ {
		ns := framework.CreateTestingNamespace(fmt.Sprintf("ns%d", i), s, t)
		defer framework.DeleteTestingNamespace(ns, s, t)
		namespaces = append(namespaces, ns)

		_, err := clientSet.CoreV1().Secrets(ns.Name).Create(makeSecret("foo"))
		if err != nil {
			t.Errorf("Couldn't create secret: %v", err)
		}
		_, err = clientSet.CoreV1().Secrets(ns.Name).Create(makeSecret("bar"))
		if err != nil {
			t.Errorf("Couldn't create secret: %v", err)
		}
	}

	testcases := []struct {
		namespace       string
		selector        string
		expectedSecrets int
	}{
		{
			namespace:       "",
			selector:        "metadata.name=foo",
			expectedSecrets: numNamespaces,
		},
		{
			namespace:       "",
			selector:        "metadata.name=foo,metadata.name=bar",
			expectedSecrets: 0,
		},
		{
			namespace:       "",
			selector:        "metadata.name=foo,metadata.namespace=ns1",
			expectedSecrets: 1,
		},
		{
			namespace:       "ns1",
			selector:        "metadata.name=foo,metadata.namespace=ns1",
			expectedSecrets: 1,
		},
		{
			namespace:       "ns1",
			selector:        "metadata.name=foo,metadata.namespace=ns2",
			expectedSecrets: 0,
		},
		{
			namespace:       "ns1",
			selector:        "metadata.name=foo,metadata.namespace=",
			expectedSecrets: 0,
		},
	}

	for _, tc := range testcases {
		opts := metav1.ListOptions{
			FieldSelector: tc.selector,
		}
		secrets, err := clientSet.CoreV1().Secrets(tc.namespace).List(opts)
		if err != nil {
			t.Errorf("%s: Unexpected error: %v", tc.selector, err)
		}
		if len(secrets.Items) != tc.expectedSecrets {
			t.Errorf("%s: Unexpected number of secrets: %d, expected: %d", tc.selector, len(secrets.Items), tc.expectedSecrets)
		}
	}
}
