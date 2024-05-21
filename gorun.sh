#!/bin/bash


# if [ "$1" = "test" ]; then
#     go test -v -race $(ls *.go | grep _test.go) *.go
# else
#     go run $(ls *.go | grep -v _test.go)
# fi



if [ "$1" = "test" ]; then
    go test -v -race .
else
    go run .
fi

# name: Go

# on:
#   push:
#     branches: [ "main" ]
#   pull_request:
#     branches: [ "main" ]

# jobs:

#   build:
#     runs-on: ubuntu-latest
#     steps:
#     - uses: actions/checkout@v4

#     - name: Set up Go
#       uses: actions/setup-go@v4
#       with:
#         go-version: '1.21.6'

#     - name: Build
#       run: go build -v ./...

#     - name: Test
#       run: go test -v ./...



