#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build command
now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GOOS=js GOARCH=wasm go build -o main.wasm $@