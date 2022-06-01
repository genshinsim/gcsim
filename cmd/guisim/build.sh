#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build command
now=$(date --utc +%FT%T%Z)
go build -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"
