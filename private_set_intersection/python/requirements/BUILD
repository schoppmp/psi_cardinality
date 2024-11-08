load("@python_versions//3.10:defs.bzl", compile_pip_requirements_3_10 = "compile_pip_requirements")
load("@python_versions//3.11:defs.bzl", compile_pip_requirements_3_11 = "compile_pip_requirements")
load("@python_versions//3.12:defs.bzl", compile_pip_requirements_3_12 = "compile_pip_requirements")
load("@python_versions//3.13:defs.bzl", compile_pip_requirements_3_13 = "compile_pip_requirements")
load("@python_versions//3.8:defs.bzl", compile_pip_requirements_3_8 = "compile_pip_requirements")
load("@python_versions//3.9:defs.bzl", compile_pip_requirements_3_9 = "compile_pip_requirements")

compile_pip_requirements_3_8(
    name = "requirements_3_8",
    src = "requirements.in",
    requirements_txt = "requirements_lock_3_8.txt",
)

compile_pip_requirements_3_9(
    name = "requirements_3_9",
    src = "requirements.in",
    requirements_txt = "requirements_lock_3_9.txt",
)

compile_pip_requirements_3_10(
    name = "requirements_3_10",
    src = "requirements.in",
    requirements_txt = "requirements_lock_3_10.txt",
)

compile_pip_requirements_3_11(
    name = "requirements_3_11",
    src = "requirements.in",
    requirements_txt = "requirements_lock_3_11.txt",
)

compile_pip_requirements_3_12(
    name = "requirements_3_12",
    src = "requirements.in",
    requirements_txt = "requirements_lock_3_12.txt",
)

compile_pip_requirements_3_13(
    name = "requirements_3_13",
    src = "requirements.in",
    requirements_txt = "requirements_lock_3_13.txt",
)
