#!/bin/bash

# TODO: get from file or env
SHARE_KEY=$GCSIM_SHARE_KEY

GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o gcsim_windows_amd64.exe ./cmd/gcsim
GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o gcsim_darwin_arm64 ./cmd/gcsim
GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o gcsim_darwin_amd64 ./cmd/gcsim
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o gcsim_linux_amd64 ./cmd/gcsim

GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}'  -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_windows_amd64.exe ./cmd/server
GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_darwin_arm64 ./cmd/server 
GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_darwin_amd64 ./cmd/server 
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.shareKey=${SHARE_KEY}' -X main.version=`git tag --sort=-version:refname | head -n 1`" -o server_linux_amd64 ./cmd/server 

go install fyne.io/fyne/v2/cmd/fyne@v2.5.2
cd cmd/server_andriod && go mod download && GOFLAGS="'-ldflags=-X=main.version=`git tag --sort=-version:refname | head -n 1` \"-X=main.shareKey=${SHARE_KEY}\"'" fyne package -os android -appID com.gcsim.server -icon ../../ui/packages/ui/src/Images/logo.png --release --name "gcsim server" && cp gcsim_server.apk ../../gcsim_server.apk

