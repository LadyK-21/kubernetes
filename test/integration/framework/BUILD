package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "etcd.go",
        "master_utils.go",
        "perf_utils.go",
        "serializer.go",
        "test_server.go",
        "util.go",
    ],
    data = [
        "@com_coreos_etcd//:etcd",
    ],
    importpath = "k8s.io/kubernetes/test/integration/framework",
    deps = [
        "//cmd/kube-apiserver/app:go_default_library",
        "//cmd/kube-apiserver/app/options:go_default_library",
        "//pkg/api/legacyscheme:go_default_library",
        "//pkg/api/testapi:go_default_library",
        "//pkg/apis/batch:go_default_library",
        "//pkg/apis/policy/v1beta1:go_default_library",
        "//pkg/generated/openapi:go_default_library",
        "//pkg/kubelet/client:go_default_library",
        "//pkg/master:go_default_library",
        "//pkg/util/env:go_default_library",
        "//pkg/version:go_default_library",
        "//staging/src/k8s.io/api/apps/v1beta1:go_default_library",
        "//staging/src/k8s.io/api/auditregistration/v1alpha1:go_default_library",
        "//staging/src/k8s.io/api/autoscaling/v1:go_default_library",
        "//staging/src/k8s.io/api/certificates/v1beta1:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/api/extensions/v1beta1:go_default_library",
        "//staging/src/k8s.io/api/rbac/v1alpha1:go_default_library",
        "//staging/src/k8s.io/api/storage/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/runtime/serializer/versioning:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/authentication/authenticator:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/authentication/authenticatorfactory:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/authentication/request/union:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/authentication/user:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/authorization/authorizer:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/authorization/authorizerfactory:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/authorization/union:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/endpoints/openapi:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/server:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/server/options:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/server/storage:go_default_library",
        "//staging/src/k8s.io/apiserver/pkg/storage/storagebackend:go_default_library",
        "//staging/src/k8s.io/client-go/informers:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes:go_default_library",
        "//staging/src/k8s.io/client-go/rest:go_default_library",
        "//staging/src/k8s.io/client-go/util/cert:go_default_library",
        "//test/e2e/framework:go_default_library",
        "//test/utils:go_default_library",
        "//vendor/github.com/go-openapi/spec:go_default_library",
        "//vendor/github.com/pborman/uuid:go_default_library",
        "//vendor/k8s.io/klog:go_default_library",
        "//vendor/k8s.io/kube-openapi/pkg/common:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
)
