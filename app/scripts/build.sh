#!/bin/bash

cd "./temp/cmd/wasm"

# notice how we avoid spaces in $now to avoid quotation hell in go build command
now=$(date +'%Y-%m-%d_%T')
GOOS=js GOARCH=wasm go build -o main.wasm -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"

mv main.wasm ../../../static