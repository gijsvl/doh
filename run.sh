#!/usr/bin/env bash

GOPATH="$(pwd)" go run src/*.go "$@"
