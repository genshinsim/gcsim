#!/bin/bash

# sudo apt-get install gcc-multilib
# sudo apt-get install gcc-mingw-w64

# notice how we avoid spaces in $now to avoid quotation hell in go build command
cd "./cmd/wasm"
now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GOOS=windows GOARCH=amd64 go build -o ../../gcsim.exe -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now" ./cmd/gcsim/ 