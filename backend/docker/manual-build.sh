#!/bin/bash

#always start at root of git repo
cd $(git rev-parse --show-cdup)

( \
cd ./backend/cmd/db \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/db \
)


( \
cd ./backend/cmd/jadechamber \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/jadechamber \
)

( \
cd ./backend/cmd/managerchan \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/managerchan \
)

( \
cd ./backend/cmd/preview \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/preview \
)

( \
cd ./backend/cmd/share \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/share \
)

( \
cd ./backend/cmd/notification \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/binary/notification \
)