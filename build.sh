#!/usr/bin/env bash

PACKAGE=github.com/jhonnyV-V/phoemux
VERSION=$1

if [[ $# -eq 0 ]]; then
  VERSION=$(git describe --abbrev=0 --tags)
fi

LDFLAGS=(
  "-s"
  "-X '${PACKAGE}/version.Version=${VERSION}'"
)

GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS[*]}" -o ./dist/phoemux_amd64 ./
GOOS=linux GOARCH=arm64 go build -ldflags="${LDFLAGS[*]}" -o ./dist/phoemux_arm64 ./

GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS[*]}" -o ./dist/phoemux_darwin_amd64 ./
GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS[*]}" -o ./dist/phoemux_darwin_arm64 ./
