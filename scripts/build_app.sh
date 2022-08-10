#!/bin/bash

cd "./cmd/wasm"
now=$(date --utc +%FT%T%Z)
GOOS=js GOARCH=wasm go build -o main.wasm -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"

mv main.wasm ../../app/static