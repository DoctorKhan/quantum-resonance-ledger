#!/usr/bin/env python3
import re
from datetime import datetime

TARGET_FILE = "docs/whitepaper.md"

def update_date():
    with open(TARGET_FILE, "r") as file:
        content = file.readlines()

    date = datetime.now().strftime("%Y-%m-%d")
    updated_content = [
        re.sub(r"^Last Updated: .*", f"Last Updated: {date}", line)
        for line in content
    ]

    with open(TARGET_FILE, "w") as file:
        file.writelines(updated_content)

if __name__ == "__main__":
    update_date()