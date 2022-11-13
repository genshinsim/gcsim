#!/bin/bash

# sudo apt-get install gcc-multilib
# sudo apt-get install gcc-mingw-w64

# notice how we avoid spaces in $now to avoid quotation hell in go build command
now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GOOS=windows GOARCH=amd64 go build ./cmd/gcsim