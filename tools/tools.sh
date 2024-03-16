#!/bin/bash

mkdir -p dist
pushd decode
  go build -o ../dist/decode
popd

pushd encode
  go build -o ../dist/encode
popd

pushd build 
  go build -o ../dist/build
popd
