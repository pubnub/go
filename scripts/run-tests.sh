#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go test -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./ ./

go test -v -coverprofile=integration_tests.out \
-covermode=atomic -coverpkg=./ ./tests/e2e/

gocovmerge functional_tests.out integration_tests.out > coverage.txt

rm integration_tests.out functional_tests.out
