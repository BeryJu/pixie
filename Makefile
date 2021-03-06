all: macos

GOPATH := $(shell go env GOPATH)

test:
	go test

lint:
	go get golang.org/x/lint/golint
	$(GOPATH)/bin/golint

macos:
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=amd64
	go build -ldflags "-s -w" -v -o bin/pixie-darwin-amd64

linux-amd64:
	go get github.com/gobuffalo/packr/v2/packr2
	$(GOPATH)/bin/packr2
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=amd64
	go build -ldflags "-s -w" -v -o bin/pixie-linux-amd64

linux-arm64:
	go get github.com/gobuffalo/packr/v2/packr2
	$(GOPATH)/bin/packr2
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=arm64
	go build -ldflags "-s -w" -v -o bin/pixie-linux-arm64

run:
	go run -v . demo --debug
