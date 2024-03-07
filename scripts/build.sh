#!/bin/bash

# TODO: get from file or env
SHARE_KEY=$GCSIM_SHARE_KEY

GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" ./cmd/gcsim
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}'  -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_windows_amd64.exe ./cmd/server
GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_darwin_arm64 ./cmd/server 
GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_darwin_amd64 ./cmd/server 
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_linux_amd64 ./cmd/server 