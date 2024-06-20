#!/bin/bash


# if [ "$1" = "test" ]; then
#     go test -v -race $(ls *.go | grep _test.go) *.go
# else
#     go run $(ls *.go | grep -v _test.go)
# fi




if [ "$1" = "test" ]; then
    cd test
    go test -v -race ./...
    cd ..
else
    go run .
fi




