#!/usr/bin/env bash

rm ./dist/*
./build.sh

echo "previous tag: $(git describe --abbrev=0 --tags)"

if [[ $# -eq 0 ]]; then
  echo "no version tag passed"
  exit
fi

gh release create $1 ./dist/* --generate-notes
