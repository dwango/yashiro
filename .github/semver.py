#!/usr/bin/env python3

from __future__ import annotations
import re
import os
import sys

# reference: https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
SEMVER_REGEX = r"^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$"


def removeprefix(text: str, prefix: str) -> str:
    if text.startswith(prefix):
        return text[len(prefix) :]
    else:
        return text


def get_semver(text: str, prefix: str) -> str | None:
    if prefix != "":
        text = removeprefix(text, prefix)

    r = re.compile(SEMVER_REGEX)
    result = r.match(text)

    if result is None:
        raise Exception(f"target is not semver string: '{text}'")

    return (
        text,
        result["major"],
        result["minor"],
        result["patch"],
        result["prerelease"],
        result["buildmetadata"],
    )


def main():
    semver_str: str = ""
    prefix: str = ""

    try:
        semver_str = sys.argv[1]
        prefix = sys.argv[2]
    except IndexError:
        pass

    output = os.getenv("GITHUB_OUTPUT")
    if output is None:
        raise Exception("GITHUB_OUTPUT environment variable is not define")

    semver, major, minor, patch, prerelease, build_metadata = get_semver(
        semver_str, prefix
    )

    if prerelease is None:
        prerelease = ""

    if build_metadata is None:
        build_metadata = ""

    with open(output, "a") as f:
        print(f"version={semver}", file=f)
        print(f"major={major}", file=f)
        print(f"minor={minor}", file=f)
        print(f"patch={patch}", file=f)
        print(f"prerelease={prerelease}", file=f)
        print(f"build-metadata={build_metadata}", file=f)


if __name__ == "__main__":
    try:
        main()
    except Exception as err:
        print(f"::error::raised error: {err}")
        sys.exit(1)
