#!/bin/bash

# TODO: get from file or env
SHARE_KEY=$GCSIM_SHARE_KEY

LDFLAGS=(
  # "-w -s" # reduces binary size at cost of performance
  "-X 'main.shareKey=${SHARE_KEY}'"
)

# reduces by ~2MB but makes really slow: -gcflags=all="-l -B -C -std"

GOOS=js GOARCH=wasm go build -o main.wasm -ldflags="${LDFLAGS[*]}" $@