#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o ./dist/phoemux_amd64 ./
GOOS=linux GOARCH=arm64 go build -ldflags='-s' -o ./dist/phoemux_arm64 ./

GOOS=darwin GOARCH=amd64 go build -ldflags='-s' -o ./dist/phoemux_darwin_amd64 ./
GOOS=darwin GOARCH=arm64 go build -ldflags='-s' -o ./dist/phoemux_darwin_arm64 ./
