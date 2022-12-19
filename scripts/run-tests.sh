#!/usr/bin/env bash
WITH_MOD=$1

go version
go env

set -e
echo "" > coverage.txt

# Set GOPATH environment if it is missing.
[[ -z "$GOPATH" ]] && GOPATH="$(go env GOPATH)"

if [[ -n "$GOPATH" ]]; then
  echo "Create module folder if required: $GOPATH/src/github.com/pubnub"
  ! [[ -d "$GOPATH/src/github.com/pubnub" ]] && mkdir -p "$GOPATH/src/github.com/pubnub"
else
  echo "'GOPATH' not defined."
fi

clean_coverage_output () {
   find . -type f -name "*.out" | xargs -r rm -rf
}

# Install 'gocovmerge' module.
go install github.com/wadey/gocovmerge@latest

echo "Run functional tests"
if ! test_run="$(go test $WITH_MOD -v -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./ ./ 2>&1)"; then
  echo "::error title=test::Functional tests failed: $test_run"
  clean_coverage_output
  exit 1
fi

echo "Run utils tests"
if ! test_run="$(go test $WITH_MOD -v -race -coverprofile=utils_tests.out -covermode=atomic -coverpkg=./ ./utils/ 2>&1)"; then
  echo "::error title=test::Unit tests failed: $test_run"
  clean_coverage_output
  exit 1
fi

echo "Run helpers tests"
if ! test_run="$(go test $WITH_MOD -v -race -coverprofile=helpers_tests.out -covermode=atomic -coverpkg=./ ./tests/helpers/ 2>&1)"; then
  echo "::error title=test::Helpers tests failed: $test_run"
  clean_coverage_output
  exit 1
fi

echo "Run integration tests"
if ! test_run="$(go test $WITH_MOD -v -race -coverprofile=integration_tests.out -covermode=atomic -coverpkg=./ ./tests/e2e/ 2>&1)"; then
  echo "::error title=test::Integration tests failed: $test_run"
  clean_coverage_output
  exit 1
fi

echo "Run deadlock tests #1"
if ! test_run="$(go test $WITH_MOD -v -race -run "TestDestroy\b" -count 20 -coverprofile=deadlock_tests.out 2>&1)"; then
  echo "::error title=test::Deadlock tests #1 failed: $test_run"
  clean_coverage_output
  exit 1
fi

echo "Run deadlock tests #2"
if ! test_run="$(go test $WITH_MOD -v -race -run "TestDestroy2\b" -count 20 -coverprofile=deadlock2_tests.out -covermode=atomic -coverpkg=./ ./tests/e2e/ 2>&1)"; then
  echo "::error title=test::Deadlock tests #2 failed: $test_run"
  clean_coverage_output
  exit 1
fi

echo "Upload coverage results"
# Send test results for analysis.
gocovmerge functional_tests.out integration_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out > coverage.txt

echo "Clean up test artifacts"
clean_coverage_output
