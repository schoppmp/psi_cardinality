import os
import posixpath
import re
import shutil
import sys

from distutils import sysconfig
import setuptools
from setuptools.command import build_ext

here = os.path.dirname(os.path.abspath(__file__))

IS_WINDOWS = sys.platform.startswith("win")


def _get_version():
    """Parse the version string from __init__.py."""
    with open(os.path.join(here, "psi_cardinality", "python", "__init__.py")) as f:
        try:
            version_line = next(line for line in f if line.startswith("__version__"))
        except StopIteration:
            raise ValueError("__version__ not defined in __init__.py")
        else:
            ns = {}
            exec(version_line, ns)  # pylint: disable=exec-used
            return ns["__version__"]


def _parse_requirements(path):
    with open(os.path.join(here, path)) as f:
        return [line.rstrip() for line in f if not (line.isspace() or line.startswith("#"))]


class BazelExtension(setuptools.Extension):
    """A C/C++ extension that is defined as a Bazel BUILD target."""

    def __init__(self, name, bazel_target):
        self.bazel_target = bazel_target
        self.relpath, self.target_name = posixpath.relpath(bazel_target, "//").split(":")
        setuptools.Extension.__init__(self, name, sources=[])


class BuildBazelExtension(build_ext.build_ext):
    """A command that runs Bazel to build a C/C++ extension."""

    def run(self):
        for ext in self.extensions:
            self.bazel_build(ext)
        build_ext.build_ext.run(self)

    def bazel_build(self, ext):
        with open("WORKSPACE", "r") as f:
            workspace_contents = f.read()

        with open("WORKSPACE", "w") as f:
            f.write(
                re.sub(
                    r'(?<=path = ").*(?=",  # May be overwritten by setup\.py\.)',
                    sysconfig.get_python_inc().replace(os.path.sep, posixpath.sep),
                    workspace_contents,
                )
            )

        if not os.path.exists(self.build_temp):
            os.makedirs(self.build_temp)

        bazel_argv = [
            "bazel",
            "build",
            ext.bazel_target,
            "--symlink_prefix=" + os.path.join(self.build_temp, "bazel-"),
            "--compilation_mode=" + ("dbg" if self.debug else "opt"),
        ]

        if IS_WINDOWS:
            # Link with python*.lib.
            for library_dir in self.library_dirs:
                bazel_argv.append("--linkopt=/LIBPATH:" + library_dir)

        self.spawn(bazel_argv)

        shared_lib_suffix = ".dll" if IS_WINDOWS else ".so"
        ext_bazel_bin_path = os.path.join(
            self.build_temp, "bazel-bin", ext.relpath, "_" + ext.target_name + shared_lib_suffix
        )
        ext_dest_path = self.get_ext_fullpath(ext.name)
        ext_dest_dir = os.path.dirname(ext_dest_path)
        if not os.path.exists(ext_dest_dir):
            os.makedirs(ext_dest_dir)
        shutil.copyfile(ext_bazel_bin_path, ext_dest_path)


setuptools.setup(
    name="psi_cardinality",
    version=_get_version(),
    description="Private Set Intersection Cardinality protocol based on ECDH and Bloom Filters.",
    keywords="privacy cryptography",
    url="https://github.com/OpenMined/PSI",
    python_requires=">=3.6",
    package_dir={"": "psi_cardinality/python"},
    packages=setuptools.find_packages("psi_cardinality/python"),
    install_requires=_parse_requirements("psi_cardinality/python/requirements.txt"),
    cmdclass=dict(build_ext=BuildBazelExtension),
    ext_modules=[BazelExtension("psi_cardinality", "//psi_cardinality/python:psi_cardinality")],
    zip_safe=False,
    classifiers=[
        "Programming Language :: Python :: 3",
        "Operating System :: OS Independent",
        "Topic :: Security :: Cryptography",
    ],
    license="Apache 2.0",
)