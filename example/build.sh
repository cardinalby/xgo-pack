#!/bin/bash

#go install github.com/cardinalby/xgo-pack@latest
#xgo-pack build

cd ../
go build .
cd example || exit 1
go mod download
../xgo-pack build