all: macos linux

macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -v -o bin/pixie-darwin-amd64

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -v -o bin/pixie-linux-amd64

run:
	go run -v . -r demo --debug
