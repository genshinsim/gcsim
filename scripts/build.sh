#!/bin/bash

# TODO: get from file or env
SHARE_KEY=$GCSIM_SHARE_KEY

LDFLAGS=(
  "-X 'main.shareKey=${SHARE_KEY}'"
)

GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS[*]}" ./cmd/gcsim