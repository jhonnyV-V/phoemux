#!/usr/bin/env bash

echo "previous tag: $(git describe --abbrev=0 --tags)"

if [[ $# -eq 0 ]]; then
  echo "no version tag passed"
  exit
fi

rm ./dist/*
./build.sh

gh release create $1 ./dist/* --generate-notes
