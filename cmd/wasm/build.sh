#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build command
# reduces by ~2MB but makes really slow: -gcflags=all="-l -B -C -std"
now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GOOS=js GOARCH=wasm go build -o main.wasm -ldflags="-w -s -d" $@