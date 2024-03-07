#!/bin/bash

#go install github.com/cardinalby/xgo-pack@latest
#xgo-pack build

cd ../
go build .
cd example
go mod download
../xgo-pack build