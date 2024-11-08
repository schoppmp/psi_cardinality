# Version 2.0.6

Feat:

- Add support for Python 3.13

# Version 2.0.5

Feat:

- Add support for Python 3.12

Fix:

- There was a bug in the pybind11 bindings involving incorrect usage for
  releasing the GIL lock which led to a segfault when using Python 3.12.
  Previous versions remain unaffected.
- Pin the python `protobuf` version to match
  [`rules_proto_grpc_python`](https://github.com/rules-proto-grpc/rules_proto_grpc/blob/master/modules/python/requirements.in#L3)
  to remove runtime warnings.

Chore:

- Removed `macos-12` from the `CD` GHA as it was building with an incorrect
  wheel name for the given OS. We're only building/publishing wheels for
  `macos-14`.
- Removed unused pip deps from `requirements.in` and updated lock files
  accordingly.

# Version 2.0.4

Chore:

- No changes, but some published python wheels may have the wrong wheel
  specifications due to a potential publishing bug in the previous release.

# Version 2.0.3

Feat:

- Add support for python 3.11
- Add support for Mac ARM64 builds

Chore:

- Use Bazel 7.3.1 with `bzlmod`
- Updated all Bazel dependencies for all languages.
- Updated the fork of private_join_and_compute (PCJ) to use Bazel 7+
  with`bzlmod` and updated all dependencies.
- Refactored CI/CD scripts to remove explicit python version in CI/CD as Bazel
  can pick up the local python within the runner. Tests on Python are hermetic
  so the local version is only important when publishing wheels.
- Minor fixes that the cpp linter was reporting and a fix in the python cpp
  bindings to use absl::MakeSpan instead of passing a std::vector.

Breaking:

- The **python** package has been renamed to `openmined_psi` to comply with the
  bazel toolchain naming conventions.

# Version 2.0.2

Chore:

- Update to use Bazel 6.4.0
- Updated some of the Bazel dependencies for Golang and Python
- Updated CI/CD to use python 3.10 for Bazel to work correctly.

# Version 2.0.1

Feat:

- The complexity of the underlying `Raw` intersection computation has improved
  from `O(nmlog(m))` -> `O(nlog(n) + max(n, m))`; however, internal protobuf
  deserialization remains as the dominant performance inhibitor for the
  `client->GetIntersection*` methods.

Fix:

- The `go` integration tests were not using the datastructure param properly.
  The fix did not result in any regression.

Chore:

- Update `C++` benchmarks to include the new `Raw` enum variant
- Misc fixes to tests which were not rounding correctly and causing CI to fail
  randomly
- Update the main README to include a description of the protocol

# Version 2.0.0

Breaking:

- The protobuf schema has changed to support additional enum variants for
  backing datastructures. This means prior protobufs using earlier library
  versions will not work with `2.0.0+`. We made this change to ensure additional
  variants will not break the schema. In general, storing serialized protobufs
  is not a pattern the protocol recommends since every intersection calculation
  should be a new protocol exchange (setup, request, response, intersection).

Feat:

- A `Raw` datastructure enum variant has been added which when specified,
  computes an intersection against 'raw' encrypted strings, but with zero
  false-positives (`fpr` is ignored). It is intended for applications that
  depend on 100% accuracy, but it may have a performance impact.

Bugfix:

- Previously, only the `CPP` bindings supported the datastructure enum to allow
  the user to select which backing datastructure to use in the protocol (`GCS`,
  `BloomFilter`, and now `Raw`). This meant we were not feature-compatible in
  the languages that depended on the (incomplete) `C` bindings. In this release,
  all languages support the same features, namely, we've updated the `C`
  bindings which allowed us to get `Go` and `Rust` to be feature compatible.

Chore:

- All languages now have common integration tests and use unit tests where
  relevant. We make a distinction between unit and integration tests. Unit tests
  are for language-specific edge cases whereas the integration tests will test
  the common APIs for correctness.

# Version 1.1.1

Chore:

- Publish new `python` wheels to support `3.8.0+`

# Version 1.1.0

Bugfix:

- Fixed an issue when computing an intersection with the client's set **larger**
  than the server's set. In this case, the intersection reported more values
  than possible. e.g. computing an intersection with a client set of 100 values
  and server with 10 values sometimes resulted in 12.

# Version 1.0.3

Bugfix:

- The Python builds were using the older `manylinux2014` tagging convention
  which was causing issues on systems that expected a specific glibc version
  (Ubuntu 20.04 uses [2.31](https://launchpad.net/ubuntu/+source/glibc)). We now
  use the latest `manylinux_x_y` tagging convention to accomodate different
  glibc versions across linux when building wheels. This means we must support
  `python 3.8.10+, 3.9.5+, 3.10.0+` and `pip >= 20.3` in accordance with [PEP
  600](https://github.com/pypa/manylinux)

# Version 1.0.2

Feat:

- This patch version relaxes the Python requirement to allow for older
  `protobuf` libraries to work with the PSI library.

# Version 1.0.1

Feat:

- Update the Python bindings to take in an optional `DataStructure` argument for
  `CreateSetupMessage`. This allows the user to customize the behavior of the
  backing datastructure - i.e. to select between the default (`GCS`) or specify
  `BLOOM_FILTER`. The previous behavior always selected `GCS` so if the
  parameter is omitted, the behavior will remain the same.

  Ex:

  ```python
  import private_set_intersection.python as psi

  c = psi.client.CreateWithNewKey(...)
  s = psi.server.CreateWithNewKey(...)

  #...

  # Defaults to GCS
  s.CreateSetupMessage(fpr, len(client_items), server_items)

  # Same as above
  s.CreateSetupMessage(fpr, len(client_items), server_items, psi.DataStructure.GCS)

  # Specify BloomFilter
  s.CreateSetupMessage(fpr, len(client_items), server_items, psi.DataStructure.BLOOM_FILTER)

  # ...
  ```

# Version 1.0.0

Breaking:

Our goal was to get CI working after 2 years of dormancy. Part of that journey
meant updating almost all dependencies. There are no functional changes in this
release; however, there were many updates to the tooling which meant we needed
to deprecate a few things.

- Python 3.6 and 3.7 are no longer supported and are superceeded by 3.8+. For M1
  macs (arm64), building python < 3.8 is problematic. There is a breaking change
  in the import structure. We simplify the importing logic to remove complexity
  between the development environment and the published wheels. See the
  [code](private_set_intersection/python/__init__.py) for an example.
- The pure Javascript builds are no longer supported or packaged. We're making
  this change because this variant is an order of magnitude slower than WASM
  (which is already slower than the native builds).
- No longer supporting rust < 2018
- No longer supporting building on Windows

Feat:

- Updated to bazel 6.0.0
- Updated to use rust 2018 edition
- Updated to support go 1.17+
- Updated to support python 3.8+
- Removed all remnants of the custom emscripten toolchain (including the git
  submodule) in favor of the native `emsdk` bazel configurations.

Chore:

- Upgraded all dev dependencies for the entire project to the latest possible
  versions that were compatable with each other. This means we have changed
  forks for PJC to another that contains the latest upstream changes but with
  updated/compatable dependencies that the project requires.
- Updated tests to use `-c opt` for all CI tests as it was noticeably faster
- Updated all dev deps for the WebAssembly builds. The WASM build configuration
  is not entirely migrated to bazel and still requires the use of a few npm
  scripts to build, compile, and publish.

# Version 0.3.5

Feat:

- Use an updated version of private_join_and_compute which uses the cleaner absl
  libraries. This also pulls in some new dependencies for PJC such as gRPC.
- Allows the python variant to be built on Windows.
- Emsdk has been bumped to the latest version (emscripten v2.0.15). The latest
  version required changes to the `cc_toolchain_config.bzl` file which addresses
  some work arounds to build successfully.

# Version 0.3.4

Feat:

- Adds native Python protobuf definitions (`python_proto_library`) instead of
  referencing the protobufs from C++. The tradeoff is a more Pythonic library
  for a small performance penalty for serialization/deserialization across the
  Python <-> C++ boundary. However, the cryptographic operations still represent
  the majority of the work in the protocol and therefore the penalty introduced
  via the extra serialization should be negligible.

# Version 0.3.3

Feat:

- Add better Python support for integration with PySyft by adding wrappers.

# Version 0.3.2

Chore:

- Update build dependencies for TypeScript.

# Version 0.3.1

Fix:

- Adds protobuf TypeScript definitions to the bundle.

# Version 0.3.0

Feat:

- Added support for Golomb Compressed Sets as a memory-efficient alternative to
  Bloom Filters.

Breaking:

- The Bloom Filter implementation now has a slightly different API to support
  bulk insertion or intersection.
- The format of the ServerSetup protocol buffer has been changed to support
  serializing both Bloom Filters and Golomb Compressed Sets.

# Version 0.2.3

Feat:

- Include TypeScript declarations inside the bundle to support TS environments
- Reduce the numbed of JS bundles by using `umd` in rollup. No more `cjs`,
  `iife` or `es` as it is not necessary.

Breaking:

- Deep import links now use a single file structure instead of nested folders

# Version 0.2.2

Feat:

- Added support for React Native with WebAssembly

Minor breaking:

- Deep import links have changed to include another level specifying the
  environment, `node`, `web`, or `worker`.

# Version 0.2.1

## JS

Fix:

- There was a misconfiguration in rollup which did not properly bundle protocol
  buffers for the browser.

Breaking:

- Deep import links have changed to now support `cjs`, `iife`, and `es`.
  Notably, `umd` has been removed and superceeded by `cjs` and `iife` due to the
  fix above. Users will need to specify the deep import paths for nodejs (`cjs`)
  and the browser (`iife` or `es`) instead of relying on the previous `umd`
  build.

# Version 0.2.0

Feat:

- Implemented protocol buffers for data-marshaling for all languages

# Version 0.1.0

Initial release!

Feat:

- Data is passed between client and server as json strings for all languages
