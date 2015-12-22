default: run

VERSION := 1.0
LDFLAGS := -X github.com/rogierlommers/poddy/internal/common.CommitHash=`git rev-parse HEAD` -X github.com/rogierlommers/poddy/internal/common.BuildDate=`date +"%d-%B-%Y/%T"`
BINARY := ./bin/poddy-${VERSION}

build:
	rice embed-go -i ./internal/poddy/
	CGO_ENABLED=0 go build -ldflags "-s $(LDFLAGS)" -a -installsuffix cgo -o ./target/poddy main.go

run:
	go run *.go
