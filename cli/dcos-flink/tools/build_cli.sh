#!/bin/bash

# exit immediately on failure
set -e

BASEDIR=$(pwd)/$(dirname "$0")
cd "$BASEDIR"

CLI_EXE_NAME=dcos-flink-service

# ---

# go
cd ..

go fmt
go get

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o $CLI_EXE_NAME".exe"
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -ldflags="-s -w" -o $CLI_EXE_NAME"-darwin"
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o $CLI_EXE_NAME"-linux"
