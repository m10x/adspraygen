#!/bin/bash
version=1.0.0

mkdir build/
rm build/*

# Windows amd64
goos=windows
goarch=amd64
GOOS=$goos GOARCH=$goarch go build -o adspraygen.exe
zip build/ADSprayGen_"$version"_"$goos"_"$goarch".zip adspraygen.exe

# Linux amd64
goos=linux
goarch=amd64
GOOS=$goos GOARCH=$goarch go build -o adspraygen
tar cfvz build/ADSprayGen_"$version"_"$goos"_"$goarch".tar.gz adspraygen

# Linux arm64
goos=linux
goarch=arm64
GOOS=$goos GOARCH=$goarch go build -o adspraygen
tar cfvz build/ADSprayGen_"$version"_"$goos"_"$goarch".tar.gz adspraygen

# Darwin/MacOS amd64
goos=darwin
goarch=amd64
GOOS=$goos GOARCH=$goarch go build -o adspraygen
tar cfvz build/ADSprayGen_"$version"_"$goos"_"$goarch".tar.gz adspraygen

# Darwin/MacOS arm64
goos=darwin
goarch=arm64
GOOS=$goos GOARCH=$goarch go build -o adspraygen
tar cfvz build/ADSprayGen_"$version"_"$goos"_"$goarch".tar.gz adspraygen

# reset
GOOS=
GOARCH=

# remove binaries
rm adspraygen
rm adspraygen.exe

# generate checksum file
find build/ -type f  \( -iname "*.tar.gz" -or -iname "*.zip" \) -exec sha256sum {} + > build/ADSprayGen_"$version"_checksums_sha256.txt