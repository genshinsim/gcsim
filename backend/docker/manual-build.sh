#!/bin/bash

#always start at root of git repo
cd $(git rev-parse --show-cdup)

( \
cd ./backend/cmd/jadechamber \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/jadechamber \
)

( \
cd ./backend/cmd/share \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/share \
)
