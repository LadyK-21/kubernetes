load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["legacy.go"],
    importmap = "k8s.io/kubernetes/vendor/k8s.io/client-go/discovery/cached",
    importpath = "k8s.io/client-go/discovery/cached",
    visibility = ["//visibility:public"],
    deps = [
        "//staging/src/k8s.io/client-go/discovery:go_default_library",
        "//staging/src/k8s.io/client-go/discovery/cached/memory:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//staging/src/k8s.io/client-go/discovery/cached/disk:all-srcs",
        "//staging/src/k8s.io/client-go/discovery/cached/memory:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)