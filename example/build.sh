#!/bin/bash

# Normally, you would install it via:
# go install github.com/cardinalby/xgo-pack@latest
# Build locally to test the current version
(cd ../ && go build -o .)
alias xgo-pack=../xgo-pack

go mod download
../xgo-pack build