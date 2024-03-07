#!/bin/bash

#go install github.com/cardinalby/xgo-pack@latest
#xgo-pack build

cd ../
go build .
cd example
../xgo-pack build