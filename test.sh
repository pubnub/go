#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(find  ./* -maxdepth 10 -type d ! -ipath '*gae*'); do
    if ls $d/*_test.go &> /dev/null; then
        go test -v -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out >> coverage.txt
            rm profile.out
        fi
    fi
done
