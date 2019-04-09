#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/src
pwd

#GOPATH="$(pwd)" go generate main
GOOS=windows go build -o binaries/$(basename "$PWD")_windows.exe
GOOS=darwin go build -o binaries/$(basename "$PWD")_macos
GOOS=linux go build -o binaries/$(basename "$PWD")_linux
