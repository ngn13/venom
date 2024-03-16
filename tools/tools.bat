@echo off

mkdir dist
pushd decode
  go build -o ../dist/decode.exe
popd

pushd encode
  go build -o ../dist/encode.exe
popd

pushd build
  go build -o ../dist/build.exe
popd build
