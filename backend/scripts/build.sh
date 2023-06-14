#!/bin/bash

( \
cd ./backend/cmd/db \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../artifacts/db \
)


( \
cd ./backend/cmd/jadechamber \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../artifacts/jadechamber \
)

( \
cd ./backend/cmd/managerchan \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../artifacts/managerchan \
)

( \
cd ./backend/cmd/preview \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../artifacts/preview \
)

( \
cd ./backend/cmd/share \
&& \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../artifacts/share \
)