#!/bin/sh
LDFLAGS="-X github.com/rogierlommers/poddy/internal/common.CommitHash=`git rev-parse HEAD` -X github.com/rogierlommers/poddy/internal/common.BuildDate=`date +"%d-%B-%Y/%T"`"
echo "start build"
echo "output: /bin/poddy"

# go get github.com/GeertJohan/go.rice
# go get github.com/GeertJohan/go.rice/rice

rice embed-go -i ./internal/poddy/
env GOOS=linux GOARCH=amd64 go build -ldflags "-s ${LDFLAGS}" -v -o ./bin/poddy main.go
