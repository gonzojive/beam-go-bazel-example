# The WORKSPACE file tells Bazel that this directory is a "workspace", which is like a project root.
# The content of this file specifies all the external dependencies Bazel needs to perform a build.
workspace(
    name = "beamgo_example",
)

# These rules are built-into Bazel but we need to load them first to download more rules
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

###################################################
# Protobuf setup
###################################################

# git_repository(
#     name = "com_google_protobuf",
#     commit = "052dc799d24a2d27d847e6bad78019364e1b188f",
#     remote = "https://github.com/protocolbuffers/protobuf.git",
# )

git_repository(
    name = "rules_proto",
    commit = "af6481970a34554c6942d993e194a9aed7987780",
    remote = "https://github.com/bazelbuild/rules_proto.git",
    shallow_since = "1610710171 +0100",
)

# Branch: master
# Commit: 7c95feba87ae269d09690fcebb18c77d8b8bcf6a
# Date: 2021-11-16 02:17:58 +0000 UTC
# URL: https://github.com/stackb/rules_proto/commit/7c95feba87ae269d09690fcebb18c77d8b8bcf6a
#
# V2 (#193)
# Size: 885598 (886 kB)
http_archive(
    name = "build_stack_rules_proto",
    sha256 = "1190c296a9f931343f70e58e5f6f9ee2331709be4e17001bb570e41237a6c497",
    strip_prefix = "rules_proto-7c95feba87ae269d09690fcebb18c77d8b8bcf6a",
    urls = ["https://github.com/stackb/rules_proto/archive/7c95feba87ae269d09690fcebb18c77d8b8bcf6a.tar.gz"],
)

register_toolchains("@build_stack_rules_proto//toolchain:standard")

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "8f5f192ba02319254aaf2cdcca00ec12eaafeb979a80a1e946773c520ae0a2c9",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/3.7.0/rules_nodejs-3.7.0.tar.gz"],
)

git_repository(
    # Commit landed on: 2020-10-30
    name = "com_google_googleapis",
    commit = "37ba54d7ed4da3052ceff96292631c3b6aae0e63",
    remote = "https://github.com/gonzojive/googleapis.git",
)

git_repository(
    name = "com_google_protobuf_cc",
    commit = "909a0f36a10075c4b4bc70fdee2c7e32dd612a72",  # release v3.17.3
    remote = "https://github.com/protocolbuffers/protobuf.git",
)

# for some reason, gazelle:repository go_repository name=build_bazel_rules_typescript importpath=github.com/bazelbuild/rules_typescript
# does not seem to work

http_archive(
    name = "bazel_toolchains",
    sha256 = "1caf8584434d3e31be674067996be787cfa511fda2a0f05811131b588886477f",
    strip_prefix = "bazel-toolchains-3.7.2",
    urls = [
        "https://github.com/bazelbuild/bazel-toolchains/releases/download/3.7.2/bazel-toolchains-3.7.2.tar.gz",
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-toolchains/archive/3.7.2.tar.gz",
    ],
)

####################################################
# C++ setup
#
# See https://github.com/bazelbuild/rules_go/issues/316 for an explanation of
# how to set up a hermetic C++ toolchain.
#
####################################################

# This sysroot is used by github.com/vsco/bazel-toolchains.
http_archive(
    name = "org_chromium_sysroot_linux_x64",
    build_file_content = """
filegroup(
  name = "sysroot",
  srcs = glob(["*/**"]),
  visibility = ["//visibility:public"],
)
""",
    sha256 = "84656a6df544ecef62169cfe3ab6e41bb4346a62d3ba2a045dc5a0a2ecea94a3",
    urls = ["https://commondatastorage.googleapis.com/chrome-linux-sysroot/toolchain/2202c161310ffde63729f29d27fe7bb24a0bc540/debian_stretch_amd64_sysroot.tar.xz"],
)

# https://github.com/grailbio/bazel-toolchain
BAZEL_TOOLCHAIN_TAG = "0.7.2"

BAZEL_TOOLCHAIN_SHA = "f7aa8e59c9d3cafde6edb372d9bd25fb4ee7293ab20b916d867cd0baaa642529"

http_archive(
    name = "com_grail_bazel_toolchain",
    canonical_id = BAZEL_TOOLCHAIN_TAG,
    sha256 = BAZEL_TOOLCHAIN_SHA,
    strip_prefix = "bazel-toolchain-{tag}".format(tag = BAZEL_TOOLCHAIN_TAG),
    url = "https://github.com/grailbio/bazel-toolchain/archive/{tag}.tar.gz".format(tag = BAZEL_TOOLCHAIN_TAG),
)

load("@com_grail_bazel_toolchain//toolchain:deps.bzl", "bazel_toolchain_dependencies")

bazel_toolchain_dependencies()

load("@com_grail_bazel_toolchain//toolchain:rules.bzl", "llvm_toolchain")

llvm_toolchain(
    name = "llvm_toolchain",
    llvm_version = "14.0.0",
    sysroot = {
        "linux-x86_64": "@org_chromium_sysroot_linux_x64//:sysroot",
    },
)

####################################################
# Go setup
####################################################

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "ab21448cef298740765f33a7f5acee0607203e4ea321219f2a4c85a6e0fb0a27",
    urls = [
        "https://github.com/bazelbuild/rules_go/releases/download/v0.32.0/rules_go-v0.32.0.zip",
        # "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.32.0/rules_go-v0.32.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "5982e5463f171da99e3bdaeff8c0f48283a7a5f396ec5282910b9e8a49c0dd7e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
    ],
)

# Load macros and repository rules.
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@llvm_toolchain//:toolchains.bzl", "llvm_register_toolchains")

llvm_register_toolchains()

# Go repos are maintained in a separate file because there are so many of them.
load("//:workspace_go_deps.bzl", "gazelle_managed_go_repositories")

# Declare indirect dependencies and register toolchains.
go_rules_dependencies()

go_register_toolchains(version = "1.18")

# gazelle:repository_macro workspace_go_deps.bzl%gazelle_managed_go_repositories
gazelle_managed_go_repositories()

gazelle_dependencies()

####################################################
# Java
#
# https://github.com/bazelbuild/rules_jvm_external
####################################################

RULES_JVM_EXTERNAL_TAG = "4.2"

RULES_JVM_EXTERNAL_SHA = "cd1a77b7b02e8e008439ca76fd34f5b07aecb8c752961f9640dea15e9e5ba1ca"

http_archive(
    name = "rules_jvm_external",
    sha256 = RULES_JVM_EXTERNAL_SHA,
    strip_prefix = "rules_jvm_external-%s" % RULES_JVM_EXTERNAL_TAG,
    url = "https://github.com/bazelbuild/rules_jvm_external/archive/%s.zip" % RULES_JVM_EXTERNAL_TAG,
)

load("@rules_jvm_external//:repositories.bzl", "rules_jvm_external_deps")

rules_jvm_external_deps()

load("@rules_jvm_external//:setup.bzl", "rules_jvm_external_setup")

rules_jvm_external_setup()

load("@rules_jvm_external//:defs.bzl", "maven_install")

maven_install(
    artifacts = [
        "org.apache.beam:beam-runners-flink-1.14:2.39.0",
        # Flink deps based on
        # https://sourcegraph.com/github.com/apache/beam/-/blob/runners/flink/flink_runner.gradle?L163:5
        "org.apache.flink:flink-clients_2.12:1.14.3",
        "org.apache.flink:flink-core:1.14.3",
        "org.apache.flink:flink-metrics-core:1.14.3",
        "org.apache.flink:flink-java:1.14.3",
        "org.apache.flink:flink-runtime:1.14.3",
        "org.apache.flink:flink-runtime-web_2.12:1.14.3",
        "org.apache.flink:flink-optimizer:1.14.3",
        "org.apache.flink:flink-streaming-java_2.12:1.14.3",
        "org.apache.flink:flink-statebackend-rocksdb_2.12:1.14.3",
        # Log4j
        "org.slf4j:slf4j-simple:1.7.36",
    ],
    fail_on_missing_checksum = False,
    fetch_sources = True,
    # https://github.com/bazelbuild/rules_jvm_external#pinning-artifacts-and-integration-with-bazels-downloader
    #
    # Command to keep the maven_install.json file up to date:
    #
    # bazel run @unpinned_maven//:pin
    maven_install_json = "//:maven_install.json",
    repositories = [
        "https://maven.google.com",
        "https://repo1.maven.org/maven2",
        "https://maven-central.storage.googleapis.com/repos/central/data/",
    ],
)

load("@maven//:defs.bzl", "pinned_maven_install")

pinned_maven_install()
