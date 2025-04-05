#!/usr/bin/env python3
import subprocess
import re
from datetime import datetime
import sys

WHITEPAPER_PATH = "docs/qrl_whitepaper.md"
DATE_PREFIX = "**Date:**"

def get_staged_files():
    """Returns a list of staged files."""
    try:
        result = subprocess.run(
            ['git', 'diff', '--cached', '--name-only'],
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout.splitlines()
    except subprocess.CalledProcessError as e:
        print(f"Error getting staged files: {e}", file=sys.stderr)
        return []
    except FileNotFoundError:
        print("Error: git command not found. Is git installed and in PATH?", file=sys.stderr)
        return []


def update_date_in_file(filepath):
    """Reads the file, updates the date line, and writes it back."""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            lines = f.readlines()
    except FileNotFoundError:
        print(f"Error: File not found {filepath}", file=sys.stderr)
        return False
    except Exception as e:
        print(f"Error reading file {filepath}: {e}", file=sys.stderr)
        return False

    updated = False
    new_lines = []
    current_date_str = datetime.now().strftime("%B %d, %Y") # e.g., April 03, 2025

    for line in lines:
        if line.strip().startswith(DATE_PREFIX):
            new_date_line = f"{DATE_PREFIX} {current_date_str}\n"
            if line != new_date_line:
                print(f"Updating date in {filepath} to {current_date_str}")
                new_lines.append(new_date_line)
                updated = True
            else:
                new_lines.append(line) # Date is already current
        else:
            new_lines.append(line)

    if updated:
        try:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.writelines(new_lines)
            # Stage the file again after modification
            subprocess.run(['git', 'add', filepath], check=True)
            print(f"Re-staged {filepath} with updated date.")
        except Exception as e:
            print(f"Error writing updated file {filepath}: {e}", file=sys.stderr)
            # Attempt to restore original content? Maybe too complex for hook.
            return False
    return True


if __name__ == "__main__":
    staged_files = get_staged_files()

    if WHITEPAPER_PATH in staged_files:
        print(f"Detected staged changes in {WHITEPAPER_PATH}. Checking date...")
        if not update_date_in_file(WHITEPAPER_PATH):
            print(f"Error updating date in {WHITEPAPER_PATH}. Aborting commit.", file=sys.stderr)
            sys.exit(1) # Exit with non-zero code to abort commit
        else:
            print("Date check/update complete.")
    else:
        # print(f"{WHITEPAPER_PATH} not staged, skipping date update.") # Optional debug msg
        pass

    sys.exit(0) # Exit with zero code to allow commit