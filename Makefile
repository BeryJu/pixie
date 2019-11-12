all: macos linux
macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -o bin/pixie-darwin-amd64

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/pixie-linux-amd64

run:
	./bin/pixie-darwin-amd64 -r demo/
