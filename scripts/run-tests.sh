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

echo "Go version: $(go version)"

# Install 'gocovmerge' module.
go get -u github.com/wadey/gocovmerge

echo "Run functional tests"
go test $WITH_MOD -v -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./ ./

echo "Run utils tests"
go test $WITH_MOD -v -race -coverprofile=utils_tests.out -covermode=atomic -coverpkg=./ ./utils/

echo "Run helpers tests"
go test $WITH_MOD -v -race -coverprofile=helpers_tests.out -covermode=atomic -coverpkg=./ ./tests/helpers/

echo "Run integration tests"
go test $WITH_MOD -v -race -coverprofile=integration_tests.out \
-covermode=atomic -coverpkg=./ ./tests/e2e/

echo "Run deadlock tests #1"
go test $WITH_MOD -v -race -run "TestDestroy\b" -count 20 -coverprofile=deadlock_tests.out

echo "Run deadlock tests #2"
go test $WITH_MOD -v -race -run "TestDestroy2\b" -count 20 -coverprofile=deadlock2_tests.out \
-covermode=atomic -coverpkg=./ ./tests/e2e/

echo "Upload coverage results"
# Send test results for analysis.
gocovmerge functional_tests.out integration_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out > coverage.txt

echo "Clean up test artifacts"
rm integration_tests.out functional_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out
