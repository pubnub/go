#!/usr/bin/env bash

# gotestsum-based test runner for PubNub Go SDK.
#
# Behavior:
#   - Runs every test step (functional, utils, helpers, integration, examples,
#     deadlock #1, deadlock #2) regardless of whether earlier steps failed, so
#     CI surfaces ALL failing suites in a single run instead of stopping at the
#     first failure.
#   - Each step's failure is recorded; the script still exits non-zero at the
#     end if any step failed, so CI correctly marks the job as failed.
#   - Coverage is merged from whichever per-step *.out files were produced
#     (Go writes coverage data even when individual tests fail, as long as the
#     package built successfully).

# NOTE: `set -e` is intentionally NOT used; we want to continue past failing
# steps. We do enable pipefail so that a failure inside a pipeline (e.g. a
# future `gotestsum | tee`) is still observable.
set -o pipefail

WITH_MOD=${1:-}
FORMAT=${2:-standard-verbose}  # Default format: standard-verbose

echo "🧪 Running tests with gotestsum (format: $FORMAT)"
go version

# Reset merged coverage file.
echo "" > coverage.txt

# Set GOPATH environment if it is missing.
[[ -z "${GOPATH:-}" ]] && GOPATH="$(go env GOPATH)"

if [[ -n "$GOPATH" ]]; then
  echo "📁 Create module folder if required: $GOPATH/src/github.com/pubnub"
  ! [[ -d "$GOPATH/src/github.com/pubnub" ]] && mkdir -p "$GOPATH/src/github.com/pubnub"
else
  echo "❌ 'GOPATH' not defined."
  go env
  exit 1
fi

# Clean up after coverage tool.
clean_coverage_output() {
  echo "🧹 Clean up test artifacts"
  find . -type f -name "*.out" | xargs -r rm -rf
}

# Common gotestsum flags (explicit for consistency across environments).
GOTESTSUM_FLAGS="--format=$FORMAT --format-hide-empty-pkg"

# Track failed steps so we can keep running and report them all at the end.
FAILED_STEPS=()

# run_step <label> <command...>
# Runs the given command, prints a banner, and on failure records the label
# in FAILED_STEPS without aborting the script.
run_step() {
  local label="$1"
  shift
  echo ""
  echo "$label"
  if ! "$@"; then
    echo "❌ ${label} FAILED (continuing with remaining steps)"
    FAILED_STEPS+=("$label")
  fi
}

echo ""
echo "📋 Test Plan (all steps run; failures are reported at the end):"
echo "  1. Functional tests (main package)"
echo "  2. Utils tests"
echo "  3. Helpers tests"
echo "  4. Integration tests"
echo "  5. Example tests (snippets/api)"
echo "  6. Deadlock tests (x20)"
echo "  7. Additional deadlock tests (x20)"
echo ""

# 1. Functional tests
run_step "🔧 Running functional tests..." \
  gotestsum $GOTESTSUM_FLAGS \
  --raw-command -- go test $WITH_MOD -json -v \
  -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./ ./

# 2. Utils tests
run_step "🛠️  Running utils tests..." \
  gotestsum $GOTESTSUM_FLAGS \
  --raw-command -- go test $WITH_MOD -json -v -race \
  -coverprofile=utils_tests.out -covermode=atomic -coverpkg=./ ./utils/

# 3. Helpers tests
run_step "🤝 Running helpers tests..." \
  gotestsum $GOTESTSUM_FLAGS \
  --raw-command -- go test $WITH_MOD -json -v -race \
  -coverprofile=helpers_tests.out -covermode=atomic -coverpkg=./ ./tests/helpers/

# 4. Integration tests
run_step "🔗 Running integration tests..." \
  gotestsum $GOTESTSUM_FLAGS \
  --raw-command -- go test $WITH_MOD -json -v -race \
  -coverprofile=integration_tests.out -covermode=atomic -coverpkg=./ ./tests/e2e/

# 5. Example tests
run_step "📚 Running example tests (snippets/api)..." \
  gotestsum $GOTESTSUM_FLAGS \
  --raw-command -- go test $WITH_MOD -json -v \
  -coverprofile=examples_tests.out -covermode=atomic -coverpkg=./ ./examples/snippets/api/

# 6. Deadlock tests #1
run_step "🔒 Running deadlock tests #1 (20 iterations)..." \
  gotestsum $GOTESTSUM_FLAGS \
  --raw-command -- go test $WITH_MOD -json -v -race \
  -run "TestDestroy\b" -count 20 -coverprofile=deadlock_tests.out

# 7. Deadlock tests #2
run_step "🔐 Running deadlock tests #2 (20 iterations)..." \
  gotestsum $GOTESTSUM_FLAGS \
  --raw-command -- go test $WITH_MOD -json -v -race \
  -run "TestDestroy2\b" -count 20 -coverprofile=deadlock2_tests.out \
  -covermode=atomic -coverpkg=./ ./tests/e2e/

# Merge whichever coverage reports were produced. Some files may be missing or
# empty if a step failed to build, so only include non-empty ones.
COVERAGE_FILES=()
for f in functional_tests.out utils_tests.out helpers_tests.out \
         integration_tests.out examples_tests.out \
         deadlock_tests.out deadlock2_tests.out; do
  [[ -s "$f" ]] && COVERAGE_FILES+=("$f")
done

if command -v gocovmerge &> /dev/null; then
  if [[ ${#COVERAGE_FILES[@]} -gt 0 ]]; then
    if [[ -n "${CODACY_PROJECT_TOKEN:-}" ]]; then
      echo "📊 Uploading coverage results..."
    else
      echo "⚠️  Code coverage not uploaded because 'CODACY_PROJECT_TOKEN' not set."
      echo "📊 Merging coverage reports locally..."
    fi
    gocovmerge "${COVERAGE_FILES[@]}" > coverage.txt || \
      echo "⚠️  gocovmerge failed; coverage.txt may be incomplete."
  else
    echo "⚠️  No coverage files produced; skipping merge."
  fi
else
  echo "⚠️  gocovmerge not available. Individual coverage files preserved."
fi

clean_coverage_output

echo ""
if [[ ${#FAILED_STEPS[@]} -eq 0 ]]; then
  echo "✅ All tests completed successfully!"
  echo "📊 Coverage report available in coverage.txt"
  exit 0
else
  echo "❌ ${#FAILED_STEPS[@]} test step(s) FAILED:"
  for s in "${FAILED_STEPS[@]}"; do
    echo "   - $s"
  done
  echo ""
  echo "📊 Partial coverage report available in coverage.txt"
  exit 1
fi
