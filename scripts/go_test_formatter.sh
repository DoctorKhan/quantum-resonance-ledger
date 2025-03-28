#!/bin/bash

# Custom Jest-like formatter for Go test output
# Colors using tput
GREEN=$(tput setaf 2)
RED=$(tput setaf 1)
YELLOW=$(tput setaf 3)
CYAN=$(tput setaf 6)
BOLD=$(tput bold)
RESET=$(tput sgr0) # Reset
INVERSE=$(tput rev)

# Start timer
START_TIME=$(date +%s.%N)

# Read test output from stdin
TEST_OUTPUT=$(cat)

# End timer
END_TIME=$(date +%s.%N)
ELAPSED_TIME=$(echo "$END_TIME - $START_TIME" | bc)
ELAPSED_TIME=$(printf "%.2fs" $ELAPSED_TIME)

# Get the package name from the test output
package_line=$(echo "$TEST_OUTPUT" | grep -E "^ok  	" | head -1)
package_name=$(echo "$package_line" | awk '{print $2}')

# Count tests
TOTAL_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^=== RUN")
PASSED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^--- PASS:")
FAILED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^--- FAIL:")
SKIPPED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^--- SKIP:")

# Print Jest-like header with PASS/FAIL indicator
if [ $FAILED_TESTS -gt 0 ]; then
  echo -e "${INVERSE} ${RED}FAIL${RESET} ${CYAN}$package_name${RESET}"
elif [ $SKIPPED_TESTS -gt 0 ]; then
  echo -e "${INVERSE} ${YELLOW}PASS${RESET} ${CYAN}$package_name${RESET}"
else
  echo -e "${INVERSE} ${GREEN}PASS${RESET} ${CYAN}$package_name${RESET}"
fi

# Process and format test output for each test
echo "$TEST_OUTPUT" | grep -E "^=== RUN" | while read -r line; do
  # Extract test name
  test_name=$(echo "$line" | sed -E 's/^=== RUN   (.*)$/\1/')
  
  # Check if this test passed or failed
  if echo "$TEST_OUTPUT" | grep -E "^--- PASS: $test_name" > /dev/null; then
    # Extract test duration
    duration=$(echo "$TEST_OUTPUT" | grep -E "^--- PASS: $test_name" | sed -E 's/.*\((.*)\)/\1/')
    echo -e "  ${GREEN}✓${RESET} $test_name ${CYAN}($duration)${RESET}"
  elif echo "$TEST_OUTPUT" | grep -E "^--- FAIL: $test_name" > /dev/null; then
    # Extract test duration
    duration=$(echo "$TEST_OUTPUT" | grep -E "^--- FAIL: $test_name" | sed -E 's/.*\((.*)\)/\1/')
    echo -e "  ${RED}✕${RESET} $test_name ${CYAN}($duration)${RESET}"
    
    # Extract failure message (this is a bit tricky with Go test output)
    failure_msg=$(echo "$TEST_OUTPUT" | grep -A 10 "^--- FAIL: $test_name" | grep -v "^--- FAIL:" | grep -v "^=== RUN" | grep -v "^--- PASS:" | head -5)
    if [ ! -z "$failure_msg" ]; then
      echo ""
      echo "    ${RED}● ${test_name}${RESET}"
      echo ""
      echo "      ${RED}Error:${RESET}"
      echo "$failure_msg" | sed 's/^/      /'
      echo ""
    fi
  elif echo "$TEST_OUTPUT" | grep -E "^--- SKIP: $test_name" > /dev/null; then
    # Extract test duration
    duration=$(echo "$TEST_OUTPUT" | grep -E "^--- SKIP: $test_name" | sed -E 's/.*\((.*)\)/\1/')
    echo -e "  ${YELLOW}○${RESET} $test_name ${CYAN}($duration)${RESET} ${YELLOW}SKIPPED${RESET}"
  fi
done

# Add a fake skipped test to demonstrate yellow color if there are no skipped tests
if [ $SKIPPED_TESTS -eq 0 ]; then
  # Add a note about skipped tests in yellow
  echo -e "\n  ${YELLOW}Note: No tests were skipped. Run with -short flag to see skipped tests.${RESET}"
fi

echo ""

# Count packages
TOTAL_PACKAGES=$(echo "$TEST_OUTPUT" | grep -c "^ok  	\|^FAIL	")
PASSED_PACKAGES=$(echo "$TEST_OUTPUT" | grep -c "^ok  	")
FAILED_PACKAGES=$((TOTAL_PACKAGES - PASSED_PACKAGES))

# Extract the total time from the test output
TOTAL_TIME=$(echo "$TEST_OUTPUT" | grep -E "^ok  	" | tail -1 | awk '{print $3}')
if [ -z "$TOTAL_TIME" ]; then
  TOTAL_TIME=$ELAPSED_TIME
fi

# Print Jest-like summary with horizontal line
echo -e "${YELLOW}$(printf '%.0s-' {1..80})${RESET}"

# Print Jest-like summary line
if [ $FAILED_TESTS -gt 0 ]; then
  echo -e "${RED}Test Suites: ${FAILED_PACKAGES} failed${RESET}, ${GREEN}${PASSED_PACKAGES} passed${RESET}, ${TOTAL_PACKAGES} total"
  echo -e "${RED}Tests:       ${FAILED_TESTS} failed${RESET}, ${GREEN}${PASSED_TESTS} passed${RESET}, ${YELLOW}${SKIPPED_TESTS} skipped${RESET}, ${TOTAL_TESTS} total"
else
  if [ $SKIPPED_TESTS -gt 0 ]; then
    echo -e "${GREEN}Test Suites: ${PASSED_PACKAGES} passed${RESET}, ${TOTAL_PACKAGES} total"
    echo -e "${GREEN}Tests:       ${PASSED_TESTS} passed${RESET}, ${YELLOW}${SKIPPED_TESTS} skipped${RESET}, ${TOTAL_TESTS} total"
  else
    echo -e "${GREEN}Test Suites: ${PASSED_PACKAGES} passed${RESET}, ${TOTAL_PACKAGES} total"
    echo -e "${GREEN}Tests:       ${PASSED_TESTS} passed${RESET}, ${TOTAL_TESTS} total"
  fi
fi

# Add snapshots line like Jest
echo -e "Snapshots:   ${YELLOW}0 total${RESET}"

# Format time output like Jest
echo -e "Time:        ${YELLOW}${TOTAL_TIME}${RESET}"

# Final line
if [ $FAILED_TESTS -gt 0 ]; then
  echo -e "${RED}Ran all test suites${RED}.${RESET}"
else
  echo -e "${YELLOW}Ran all test suites${GREEN}.${RESET}"
fi