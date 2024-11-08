# PSI - Python

Private Set Intersection protocol based on ECDH, Bloom Filters, and Golomb
Compressed Sets.

## Installation

```bash
pip install openmined-psi
```

## Developing

### Tests

```bash
bazel test --test_output=all //private_set_intersection/python:test_3_8
bazel test --test_output=all //private_set_intersection/python:test_3_9
bazel test --test_output=all //private_set_intersection/python:test_3_10
bazel test --test_output=all //private_set_intersection/python:test_3_11
bazel test --test_output=all //private_set_intersection/python:test_3_12
bazel test --test_output=all //private_set_intersection/python:test_3_13
...
```

### Benchmarks

```bash
bazel run -c opt --test_output=all //private_set_intersection/python:benchmark_3_8
bazel run -c opt --test_output=all //private_set_intersection/python:benchmark_3_9
bazel run -c opt --test_output=all //private_set_intersection/python:benchmark_3_10
bazel run -c opt --test_output=all //private_set_intersection/python:benchmark_3_11
bazel run -c opt --test_output=all //private_set_intersection/python:benchmark_3_12
bazel run -c opt --test_output=all //private_set_intersection/python:benchmark_3_13
...
```

### Updating dependencies

Add any dependencies to `requirements.in`, then run:

```bash
bazel run //private_set_intersection/python/requirements:requirements_3_8.update
bazel run //private_set_intersection/python/requirements:requirements_3_9.update
bazel run //private_set_intersection/python/requirements:requirements_3_10.update
bazel run //private_set_intersection/python/requirements:requirements_3_11.update
bazel run //private_set_intersection/python/requirements:requirements_3_12.update
bazel run //private_set_intersection/python/requirements:requirements_3_13.update
...
```

### Publishing

Build the wheel

```bash
bazel build -c opt //private_set_intersection/python:wheel
```

Or build and publish in one go to test PyPi:

```bash
bazel run -c opt //private_set_intersection/python:wheel.publish  -- --repository testpypi --verbose --skip-existing
```
