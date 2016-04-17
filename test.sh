#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go version

if [[ $TRAVIS_GO_VERSION == 1.4.3 ]]; then
  echo "match"
  go get golang.org/x/tools/cmd/cover
fi

go test -coverprofile=unit_tests.out -covermode=atomic -coverpkg=./messaging ./messaging/
go test -coverprofile=functional_tests.out -covermode=atomic -coverpkg=./messaging ./messaging/tests/

gocovmerge unit_tests.out functional_tests.out > coverage.txt

rm unit_tests.out functional_tests.out
