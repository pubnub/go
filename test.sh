#!/usr/bin/env bash

set -e
echo "" > coverage.txt
go version

go test -coverprofile=unit_tests.out -covermode=atomic -coverpkg=./messaging ./messaging/
go test -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./messaging ./messaging/tests/

gocovmerge unit_tests.out functional_tests.out > coverage.txt

rm unit_tests.out functional_tests.out
