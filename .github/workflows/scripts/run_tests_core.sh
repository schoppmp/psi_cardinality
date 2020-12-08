#!/bin/sh
set -e

bazel clean --expunge

# C++
bazel test --test_output=all //private_set_intersection/cpp/...

# C
bazel test --test_output=all //private_set_intersection/c/...
