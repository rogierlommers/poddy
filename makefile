default: run

VERSION := 1.0
LDFLAGS := -X github.com/rogierlommers/poddy/internal/common.CommitHash=`git rev-parse HEAD` -X github.com/rogierlommers/poddy/internal/common.BuildDate=`date +"%d-%B-%Y/%T"`
BINARY := ./bin/poddy-${VERSION}

setup:
	go get github.com/tools/godep

build: setup
	rm -rf ./target
	mkdir -p ./target
	CGO_ENABLED=0 godep go build -ldflags "-s $(LDFLAGS)" -a -installsuffix cgo -o ./target/poddy main.go

run:
	godep go run *.go

validate:
	golint ./...
	go vet ./...

release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 godep go build -ldflags "-s $(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-darwin-386 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 godep go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-darwin-amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 godep go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-linux-386 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 godep go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-linux-amd64 main.go
	zip -m -9 $(BINARY)-darwin-386.zip $(BINARY)-darwin-386
	zip -m -9 $(BINARY)-darwin-amd64.zip $(BINARY)-darwin-amd64
	zip -m -9 $(BINARY)-linux-386.zip $(BINARY)-linux-386
	zip -m -9 $(BINARY)-linux-amd64.zip $(BINARY)-linux-amd64
