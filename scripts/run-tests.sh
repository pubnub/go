#!/usr/bin/env bash

# Colors for console.
INFO_BG="\033[0m\033[48;2;5;49;70m"
INFO_FG="\033[38;2;19;181;255m"
BOLD_INFO_FG="${INFO_FG}\033[1m"
SUCCESS_BG="\033[0m\033[48;2;30;69;1m"
SUCCESS_FG="\033[38;2;95;215;0m"
BOLD_SUCCESS_FG="${SUCCESS_FG}\033[1m"
ERROR_BG="\033[0m\033[48;2;170;86;88m"
ERROR_FG="\033[38;2;243;123;125m"
BOLD_ERROR_FG="${ERROR_FG}\033[1m"
CLEAR="\033[0m"

WITH_MOD=$1

go version
go env

set -e
echo "" > coverage.txt

# Set GOPATH environment if it is missing.
[[ -z "$GOPATH" ]] && GOPATH="$(go env GOPATH)"

if [[ -n "$GOPATH" ]]; then
  echo "::notice title=setup::Create module folder if required: $GOPATH/src/github.com/pubnub"
  ! [[ -d "$GOPATH/src/github.com/pubnub" ]] && mkdir -p "$GOPATH/src/github.com/pubnub"
else
  echo "::error title=coverage::'GOPATH' not defined."
  exit 1
fi

# Parse test suite output.
parse_test_utput () {
  TES_SUIITE_BASE="github.com/pubnub/go/v7/"
  TEST_RESULTS=""
  TEST_SUITE_RESULTS=""
  TEST_CASE_NAME=""
  TEST_CASE_BODY=""
  TEST_CASE_STARTED=0
  TEST_NAME=""
  TEST_ERROR=""
  TEST_ASSERT_LINE=""

  while IFS= read -r line; do
    # Check whether should start processing test case result.
    if [[ $TEST_CASE_STARTED -eq 0 && $line =~ ^===[[:space:]]RUN[[:space:]]+(.*)$ ]]; then
      TEST_CASE_NAME="${BASH_REMATCH[1]}"
      TEST_CASE_STARTED=1
      # Skip to next line.
      continue
    fi

    # Check whether should end processing test case result.
    if [[ $TEST_CASE_STARTED -eq 1 && $line =~ ^---[[:space:]](FAIL|PASS):[[:space:]]+(.*)[[:space:]]\((.*)\)$ ]]; then
      if [[ ${BASH_REMATCH[1]} == "PASS" ]]; then
        FORMATTED_OUTPUT="    ${BOLD_SUCCESS_FG}✓${CLEAR} $TEST_CASE_NAME (${INFO_FG}${BASH_REMATCH[3]}${CLEAR})${CLEAR}"
      else
        FORMATTED_OUTPUT="    ${BOLD_ERROR_FG}✗${CLEAR} $TEST_CASE_NAME (${INFO_FG}${BASH_REMATCH[3]}${CLEAR})${CLEAR}"
        FORMATTED_OUTPUT="$FORMATTED_OUTPUT\n$TEST_CASE_BODY"
      fi
      [[ -z $TEST_SUITE_RESULTS ]] && TEST_SUITE_RESULTS="$FORMATTED_OUTPUT" || TEST_SUITE_RESULTS="$TEST_SUITE_RESULTS\n$FORMATTED_OUTPUT"
      TEST_CASE_STARTED=0
      TEST_CASE_BODY=""
      continue
    fi
    
    if [[ -n $TEST_NAME && -n $TEST_ERROR ]]; then
      FORMATTED_ERROR="        ${BOLD_ERROR_FG}Test:${CLEAR}  ${ERROR_FG}$TEST_NAME:$TEST_ASSERT_LINE${CLEAR}\n        ${BOLD_ERROR_FG}Error:${CLEAR} ${ERROR_FG}$TEST_ERROR${CLEAR}\n"
      [[ -z $TEST_CASE_BODY ]] && TEST_CASE_BODY="$FORMATTED_ERROR" || TEST_CASE_BODY="$TEST_CASE_BODY\n$FORMATTED_ERROR"
      TEST_NAME=""
      TEST_ERROR=""
      TEST_ASSERT_LINE=""
      continue
    elif [[ -z $TEST_ERROR && $line =~ ^[[:space:]]+Error[[:space:]]Trace:[[:space:]]+.*:([0-9]+)$ ]]; then
      TEST_ASSERT_LINE="${BASH_REMATCH[1]}"
      continue
    elif [[ -z $TEST_NAME && $line =~ ^[[:space:]]+Test:[[:space:]]+(.*)$ ]]; then
      TEST_NAME="${BASH_REMATCH[1]}"
      continue
    elif [[ -z $TEST_ERROR && $line =~ ^[[:space:]]+Error:[[:space:]]+(.*)$ ]]; then
      TEST_ERROR="${BASH_REMATCH[1]}"
      continue
    fi
    
    # Check whether test suite completed or not.
    if [[ -n $TEST_SUITE_RESULTS && $line =~ ^(ok|FAIL)[[:space:]]+(.*)[[:space:]]+([0-9.]+s) ]]; then
      TEST_SUITE_NAME="${BASH_REMATCH[2]/$TES_SUIITE_BASE/}"
      if [[ ${BASH_REMATCH[1]} == "ok" ]]; then
        FORMATTED_OUTPUT="  ${BOLD_SUCCESS_FG}✓ $TEST_SUITE_NAME${CLEAR} (${INFO_FG}${BASH_REMATCH[3]}${CLEAR})${CLEAR}"
      else
        FORMATTED_OUTPUT="  ${BOLD_ERROR_FG}✗ $TEST_SUITE_NAME${CLEAR} (${INFO_FG}${BASH_REMATCH[3]}${CLEAR})${CLEAR}"
      fi
      FORMATTED_OUTPUT="$FORMATTED_OUTPUT\n$TEST_SUITE_RESULTS\n"
      
      [[ -z $TEST_RESULTS ]] && TEST_RESULTS="$FORMATTED_OUTPUT" || TEST_RESULTS="$TEST_RESULTS\n$FORMATTED_OUTPUT"
    fi
  done <<< "$0"

  echo "$TEST_RESULTS"
}

# Clean up after codecoverage tool.
clean_coverage_output () {
  echo "::notice title=coverage::Clean up test artifacts"
  find . -type f -name "*.out" | xargs -r rm -rf
}


# Install 'gocovmerge' module.
go install github.com/wadey/gocovmerge@latest

TESTS_RUN_OUTPUT=""

echo "::notice title=test::Run functional tests"
if ! TESTS_RUN_OUTPUT="$(go test $WITH_MOD -v -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./ ./ 2>&1)"; then
  echo "::error title=test::Functional tests failed: $TESTS_RUN_OUTPUT"
  clean_coverage_output
  exit 1
fi

[[ -n $TESTS_RUN_OUTPUT ]] && echo "THERE IS DATA TO PARSE" || echo "NOTHING TO PARSE"
[[ -n $TESTS_RUN_OUTPUT ]] && echo -e "$(parse_test_utput "$TESTS_RUN_OUTPUT")"
TESTS_RUN_OUTPUT=""

echo "::notice title=test::Run utils tests"
if ! TESTS_RUN_OUTPUT="$(go test $WITH_MOD -v -race -coverprofile=utils_tests.out -covermode=atomic -coverpkg=./ ./utils/ 2>&1)"; then
  echo "::error title=test::Unit tests failed: $TESTS_RUN_OUTPUT"
  clean_coverage_output
  exit 1
fi

[[ -n $TESTS_RUN_OUTPUT ]] && echo "THERE IS DATA TO PARSE" || echo "NOTHING TO PARSE"
[[ -n $TESTS_RUN_OUTPUT ]] && echo -e "$(parse_test_utput "$TESTS_RUN_OUTPUT")"
TESTS_RUN_OUTPUT=""

echo "::notice title=test::Run helpers tests"
if ! TESTS_RUN_OUTPUT="$(go test $WITH_MOD -v -race -coverprofile=helpers_tests.out -covermode=atomic -coverpkg=./ ./tests/helpers/ 2>&1)"; then
  echo "::error title=test::Helpers tests failed: $TESTS_RUN_OUTPUT"
  clean_coverage_output
  exit 1
fi

[[ -n $TESTS_RUN_OUTPUT ]] && echo "THERE IS DATA TO PARSE" || echo "NOTHING TO PARSE"
[[ -n $TESTS_RUN_OUTPUT ]] && echo -e "$(parse_test_utput "$TESTS_RUN_OUTPUT")"
TESTS_RUN_OUTPUT=""

echo "::notice title=test::Run integration tests"
if ! TESTS_RUN_OUTPUT="$(go test $WITH_MOD -v -race -coverprofile=integration_tests.out -covermode=atomic -coverpkg=./ ./tests/e2e/ 2>&1)"; then
  echo "::error title=test::Integration tests failed: $TESTS_RUN_OUTPUT"
  clean_coverage_output
  exit 1
fi

[[ -n $TESTS_RUN_OUTPUT ]] && echo "THERE IS DATA TO PARSE" || echo "NOTHING TO PARSE"
[[ -n $TESTS_RUN_OUTPUT ]] && echo -e "$(parse_test_utput "$TESTS_RUN_OUTPUT")"
TESTS_RUN_OUTPUT=""

echo "::notice title=test::Run deadlock tests #1"
if ! TESTS_RUN_OUTPUT="$(go test $WITH_MOD -v -race -run "TestDestroy\b" -count 20 -coverprofile=deadlock_tests.out 2>&1)"; then
  echo "::error title=test::Deadlock tests #1 failed: $TESTS_RUN_OUTPUT"
  clean_coverage_output
  exit 1
fi

[[ -n $TESTS_RUN_OUTPUT ]] && echo "THERE IS DATA TO PARSE" || echo "NOTHING TO PARSE"
[[ -n $TESTS_RUN_OUTPUT ]] && echo -e "$(parse_test_utput "$TESTS_RUN_OUTPUT")"
TESTS_RUN_OUTPUT=""

echo "::notice title=test::Run deadlock tests #2"
if ! TESTS_RUN_OUTPUT="$(go test $WITH_MOD -v -race -run "TestDestroy2\b" -count 20 -coverprofile=deadlock2_tests.out -covermode=atomic -coverpkg=./ ./tests/e2e/ 2>&1)"; then
  echo "::error title=test::Deadlock tests #2 failed: $TESTS_RUN_OUTPUT"
  clean_coverage_output
  exit 1
fi

[[ -n $TESTS_RUN_OUTPUT ]] && echo "THERE IS DATA TO PARSE" || echo "NOTHING TO PARSE"
[[ -n $TESTS_RUN_OUTPUT ]] && echo -e "$(parse_test_utput "$TESTS_RUN_OUTPUT")"
TESTS_RUN_OUTPUT=""

if [[ -n "$CODACY_PROJECT_TOKEN" ]]; then
  echo "::notice title=coverage::Upload coverage results"
  # Send test results for analysis.
  gocovmerge functional_tests.out integration_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out > coverage.txt
else
  echo "::warning title=coverage::Code coverage not uploaded because 'CODACY_PROJECT_TOKEN' not set."
fi

clean_coverage_output
