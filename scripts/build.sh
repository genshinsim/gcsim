#!/bin/bash
sudo apt-get install gcc-multilib
sudo apt-get install gcc-mingw-w64
# notice how we avoid spaces in $now to avoid quotation hell in go build command
now=$(date --utc +%FT%T%Z)
GOOS=windows GOARCH=amd64 go build -o gcsim.exe -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now" ./cmd/gcsim/ 
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o guisim.exe -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now" ./cmd/guisim/ 

# ls