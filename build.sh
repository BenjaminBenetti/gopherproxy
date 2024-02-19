#!/bin/bash

go build -a -v -o ./bin/gopherproxy ./cmd/gopherproxyclient
go build -a -v -o ./bin/gopherproxyserver ./cmd/gopherproxyserver