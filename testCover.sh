#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./... ); do
    if [[ "$d" == */godu ]]; then
        go test $d
    else
        go test -v -race -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out >> coverage.txt
            rm profile.out
        fi
    fi
done

go test  .
