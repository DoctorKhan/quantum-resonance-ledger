#!/usr/bin/env python3

import re

TARGET_FILE = "docs/whitepaper.md"

def increment_version():
    with open(TARGET_FILE, "r") as file:
        content = file.readlines()

    updated_content = []
    for line in content:
        if line.startswith("Version:"):
            match = re.search(r"(\d+)\.(\d+)\.(\d+)", line)
            if match:
                major, minor, patch = map(int, match.groups())
                patch += 1
                new_version = f"{major}.{minor}.{patch}"
                line = re.sub(r"\d+\.\d+\.\d+", new_version, line)
        updated_content.append(line)

    with open(TARGET_FILE, "w") as file:
        file.writelines(updated_content)

if __name__ == "__main__":
    increment_version()
