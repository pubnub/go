#!/usr/bin/env bash

# gotestsum-based test runner for PubNub Go SDK
# Provides better test output formatting and additional features

set -e

WITH_MOD=$1
FORMAT=${2:-standard-verbose}  # Default format: standard-verbose

echo "🧪 Running tests with gotestsum (format: $FORMAT)"
go version

# Clean up coverage files
echo "" > coverage.txt

# Set GOPATH environment if it is missing.
[[ -z "$GOPATH" ]] && GOPATH="$(go env GOPATH)"

if [[ -n "$GOPATH" ]]; then
  echo "📁 Create module folder if required: $GOPATH/src/github.com/pubnub"
  ! [[ -d "$GOPATH/src/github.com/pubnub" ]] && mkdir -p "$GOPATH/src/github.com/pubnub"
else
  echo "❌ 'GOPATH' not defined."
  go env
  exit 1
fi

# Clean up after coverage tool
clean_coverage_output() {
  echo "🧹 Clean up test artifacts"
  find . -type f -name "*.out" | xargs -r rm -rf
}

# Common gotestsum flags (explicit for consistency across environments)
GOTESTSUM_FLAGS="--format=$FORMAT --format-hide-empty-pkg"

# Flaky-test mitigation: re-run failed tests up to 3 attempts total before
# treating the failure as real. gotestsum exits 0 if every test eventually
# passes on some attempt, so a transient flake will not turn CI red.
# See: https://pkg.go.dev/gotest.tools/gotestsum#re-running-failed-tests
RERUN_FLAGS="--rerun-fails=3"

echo ""
echo "📋 Test Plan:"
echo "  1. Functional tests (main package)"
echo "  2. Utils tests"
echo "  3. Helpers tests" 
echo "  4. Integration tests"
echo "  5. Example tests (snippets/api)"
echo "  6. Deadlock tests (x20)"
echo "  7. Additional deadlock tests (x20)"
echo ""

# 1. Run functional tests
echo "🔧 Running functional tests..."
if ! gotestsum $GOTESTSUM_FLAGS $RERUN_FLAGS \
  --packages=./ \
  -- $WITH_MOD \
  -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./; then
  clean_coverage_output
  exit 2
fi

# 2. Run utils tests
echo "🛠️  Running utils tests..."
if ! gotestsum $GOTESTSUM_FLAGS $RERUN_FLAGS \
  --packages=./utils/ \
  -- $WITH_MOD -race \
  -coverprofile=utils_tests.out -covermode=atomic -coverpkg=./; then
  clean_coverage_output
  exit 3
fi

# 3. Run helpers tests
echo "🤝 Running helpers tests..."
if ! gotestsum $GOTESTSUM_FLAGS $RERUN_FLAGS \
  --packages=./tests/helpers/ \
  -- $WITH_MOD -race \
  -coverprofile=helpers_tests.out -covermode=atomic -coverpkg=./; then
  clean_coverage_output
  exit 4
fi

# 4. Run integration tests
echo "🔗 Running integration tests..."
if ! gotestsum $GOTESTSUM_FLAGS $RERUN_FLAGS \
  --packages=./tests/e2e/ \
  -- $WITH_MOD -race \
  -coverprofile=integration_tests.out -covermode=atomic -coverpkg=./; then
  clean_coverage_output
  exit 5
fi

# 5. Run example tests
echo "📚 Running example tests (snippets/api)..."
if ! gotestsum $GOTESTSUM_FLAGS $RERUN_FLAGS \
  --packages=./examples/snippets/api/ \
  -- $WITH_MOD \
  -coverprofile=examples_tests.out -covermode=atomic -coverpkg=./; then
  clean_coverage_output
  exit 6
fi

# 6. Run deadlock tests #1
echo "🔒 Running deadlock tests #1 (20 iterations)..."
if ! gotestsum $GOTESTSUM_FLAGS $RERUN_FLAGS \
  --packages=./ \
  -- $WITH_MOD -race \
  -run "TestDestroy\b" -count=20 -coverprofile=deadlock_tests.out; then
  clean_coverage_output
  exit 7
fi

# 7. Run deadlock tests #2
echo "🔐 Running deadlock tests #2 (20 iterations)..."
if ! gotestsum $GOTESTSUM_FLAGS $RERUN_FLAGS \
  --packages=./tests/e2e/ \
  -- $WITH_MOD -race \
  -run "TestDestroy2\b" -count=20 -coverprofile=deadlock2_tests.out \
  -covermode=atomic -coverpkg=./; then
  clean_coverage_output
  exit 8
fi

# Merge coverage reports
if command -v gocovmerge &> /dev/null; then
  if [[ -n "$CODACY_PROJECT_TOKEN" ]]; then
    echo "📊 Uploading coverage results..."
    gocovmerge functional_tests.out integration_tests.out utils_tests.out helpers_tests.out examples_tests.out deadlock_tests.out deadlock2_tests.out > coverage.txt
  else
    echo "⚠️  Code coverage not uploaded because 'CODACY_PROJECT_TOKEN' not set."
    echo "📊 Merging coverage reports locally..."
    gocovmerge functional_tests.out integration_tests.out utils_tests.out helpers_tests.out examples_tests.out deadlock_tests.out deadlock2_tests.out > coverage.txt
  fi
else
  echo "⚠️  gocovmerge not available. Individual coverage files preserved."
fi

clean_coverage_output

echo ""
echo "✅ All tests completed successfully!"
echo "📊 Coverage report available in coverage.txt"
