#!/usr/bin/env bash
go version
go env

set -e
echo "" > coverage.txt

go test -v -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./ ./

go test -v -race -coverprofile=utils_tests.out -covermode=atomic -coverpkg=./ ./utils/

go test -v -race -coverprofile=helpers_tests.out -covermode=atomic -coverpkg=./ ./tests/helpers/

go test -v -race -coverprofile=integration_tests.out \
-covermode=atomic -coverpkg=./ ./tests/e2e/

go test -v -race -run "TestDestroy\b" -count 20 -coverprofile=deadlock_tests.out \

go test -v -race -run "TestDestroy2\b" -count 20 -coverprofile=deadlock2_tests.out \
-covermode=atomic -coverpkg=./ ./tests/e2e/

gocovmerge functional_tests.out integration_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out > coverage.txt

rm integration_tests.out functional_tests.out utils_tests.out helpers_tests.out deadlock_tests.out deadlock2_tests.out
