abi_bzl_template = """\
def python_abi():
    return "{python_abi}"
"""

platform_bzl_template = """\
def python_platform():
    return "{python_platform}"
"""

# we reuse the PYTHON_BIN_PATH environment variable from pybind11 so that the
# ABI tag we detect is always compatible with the version of python that was
# used for the build
_PYTHON_BIN_PATH_ENV = "PYTHON_BIN_PATH"

def _get_python_bin(rctx):
    python_version = rctx.attr.python_version
    if python_version != "2" and python_version != "3":
        fail("python_version must be one of: '2', '3'")

    python_bin = rctx.os.environ.get(_PYTHON_BIN_PATH_ENV)
    if python_bin != None:
        return python_bin

    python_bin = rctx.which("python" + python_version)
    if python_bin != None:
        return python_bin

    fail("cannot find python binary")

def _get_python_abi(rctx, python_bin):
    result = rctx.execute([python_bin, "-c", """
import platform
assert platform.python_implementation() == "CPython"
version = platform.python_version_tuple()
print(f"cp{version[0]}{version[1]}")
"""])
    return result.stdout.splitlines()[0]

def _get_python_platform(rctx, python_bin):
    result = rctx.execute([python_bin, "-c", """
import platform
from enum import Enum
from platform import python_version

from packaging import version


class ABIVersion(Enum):
    MANY_LINUX_X_Y = "manylinux_" + str(platform.libc_ver()[1].replace(".", "_"))
    MANY_LINUX_2014 = "manylinux2014"
    MANY_LINUX_2010 = "manylinux2010"


# We rely on the bundled pip versions installed with python to determine the
# naming variant: https://github.com/pypa/manylinux
#
# NOTE: Order matters
python_versions = {
    version.parse("3.13.0"): ABIVersion.MANY_LINUX_X_Y,
    version.parse("3.12.0"): ABIVersion.MANY_LINUX_X_Y,
    version.parse("3.11.0"): ABIVersion.MANY_LINUX_X_Y,
    version.parse("3.10.0"): ABIVersion.MANY_LINUX_X_Y,
    version.parse("3.9.5"): ABIVersion.MANY_LINUX_X_Y,
    version.parse("3.9.0"): ABIVersion.MANY_LINUX_2014,
    version.parse("3.8.10"): ABIVersion.MANY_LINUX_X_Y,
    version.parse("3.8.4"): ABIVersion.MANY_LINUX_2014,
    version.parse("3.8.0"): ABIVersion.MANY_LINUX_2010,  # Fails for aarch64
}


def main():
    system = platform.system()

    manylinux = "none"
    if system.lower() == "linux":
        current_version = version.parse(python_version())
        manylinux = ABIVersion.MANY_LINUX_X_Y.value
        for ver, ml in python_versions.items():
            if current_version >= ver:
                manylinux = ml.value
                break

    print(manylinux)


if __name__ == "__main__":
    main()

"""])
    return result.stdout.splitlines()[0]

def _make_abi_impl(rctx):
    python_bin = _get_python_bin(rctx)
    python_abi = _get_python_abi(rctx, python_bin)
    rctx.file("BUILD")
    rctx.file("abi.bzl", abi_bzl_template.format(python_abi = python_abi))

def _make_platform_impl(rctx):
    python_bin = _get_python_bin(rctx)
    python_platform = _get_python_platform(rctx, python_bin)
    rctx.file("BUILD")
    rctx.file("platform.bzl", platform_bzl_template.format(python_platform = python_platform))

abi = repository_rule(
    doc = """Make variable tags to pass to a py_wheel rule.""",
    implementation = _make_abi_impl,
    attrs = {
        "python_version": attr.string(
            doc = "Python version for the build '2' or '3' accepted.",
            mandatory = True,
        ),
    },
)

platform = repository_rule(
    doc = """Make variable tags to pass to a py_wheel rule.""",
    implementation = _make_platform_impl,
    attrs = {
        "python_version": attr.string(
            doc = "Python version for the build '2' or '3' accepted.",
            mandatory = True,
        ),
    },
)

def _abi_ext_impl(_rctx):
    abi(name = "python_abi", python_version = "3")

def _platform_ext_impl(_rctx):
    platform(name = "python_platform", python_version = "3")

python_abi = module_extension(
    implementation = _abi_ext_impl,
)

python_platform = module_extension(
    implementation = _platform_ext_impl,
)
