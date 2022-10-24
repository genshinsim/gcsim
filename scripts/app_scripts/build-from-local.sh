#!/bin/bash

cd "../gcsim/cmd/wasm"

# notice how we avoid spaces in $now to avoid quotation hell in go build command
now=$(date --utc +%FT%T%Z)
GOOS=js GOARCH=wasm go build -o main.wasm -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"

mv main.wasm ../../../gcsim.app/static