#!/bin/bash

set -e

TAGS=()
if [ -n "$GCSIM_SHARE_KEY" ]; then
  go run ../sharekeygen
  TAGS=(-tags sharekey)
fi

# reduces by ~2MB but makes really slow: -gcflags=all="-l -B -C -std"
GOOS=js GOARCH=wasm go build -trimpath "${TAGS[@]}" -o main.wasm $@
