#!/usr/bin/env bash
WITH_MOD=$1

go version
go env

set -e
echo "" > coverage.txt

[[ -n "$GOPATH" && ! -d "$GOPATH/src/github.com/pubnub" ]] && \
  mkdir -p "$GOPATH/src/github.com/pubnub"

# Install 'gocovmerge' module.
go get -u github.com/wadey/gocovmerge

go test $WITH_MOD -v -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./ ./

go test $WITH_MOD -v -race -coverprofile=utils_tests.out -covermode=atomic -coverpkg=./ ./utils/

go test $WITH_MOD -v -race -coverprofile=helpers_tests.out -covermode=atomic -coverpkg=./ ./tests/helpers/

go test $WITH_MOD -v -race -coverprofile=integration_tests.out \
-covermode=atomic -coverpkg=./ ./tests/e2e/

go test $WITH_MOD -v -race -run "TestDestroy\b" -count 20 -coverprofile=deadlock_tests.out \

go test $WITH_MOD -v -race -run "TestDestroy2\b" -count 20 -coverprofile=deadlock2_tests.out \
-covermode=atomic -coverpkg=./ ./tests/e2e/

# Send test results for analysis.
gocovmerge functional_tests.out integration_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out > coverage.txt

rm integration_tests.out functional_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out
