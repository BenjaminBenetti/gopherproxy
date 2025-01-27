#!/bin/bash
if [ -z "$1" ]; then
  echo "Error: Tag argument is required"
  exit 1
fi
TAG="$1"

# build x86
go build -a -v -o ./bin/gopherproxy-${TAG}.x86 ./cmd/gopherproxyclient
go build -a -v -o ./bin/gopherproxyserver-${TAG}.x86 ./cmd/gopherproxyserver

# build arm
GOARCH=arm64 go build -a -v -o ./bin/gopherproxy-${TAG}.arm ./cmd/gopherproxyclient
GOARCH=arm64 go build -a -v -o ./bin/gopherproxyserver-${TAG}.arm ./cmd/gopherproxyserver