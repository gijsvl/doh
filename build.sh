#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

GOPATH="$(pwd)" go generate main 
GOOS=windows go build -o binaries/$(basename "$PWD")_windows.exe main
GOOS=darwin go build -o binaries/$(basename "$PWD")_macos main
GOOS=linux go build -o binaries/$(basename "$PWD")_linux main
