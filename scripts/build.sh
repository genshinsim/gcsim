#!/bin/bash

cd "./cmd/wasm"
now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GOOS=js GOARCH=wasm go build -o ../../app/static/main.wasm -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"
