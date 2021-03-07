#!/bin/sh
set -e

# Rust
# Extra flags needed for bazel 4.0.0: --incompatible_require_linker_input_cc_api=false --features=-supports_dynamic_linker
bazel test --test_output=all --incompatible_require_linker_input_cc_api=false --features=-supports_dynamic_linker //private_set_intersection/rust/...
