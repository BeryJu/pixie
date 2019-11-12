FROM golang:latest AS builder
WORKDIR $GOPATH/src/git.beryju.org/BeryJu.org/pixie
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go get github.com/gobuffalo/packr/v2/packr2 && \
    packr2 && \
    go build -v -o /go/bin/pixie

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/bin/pixie /pixie
EXPOSE 8080
CMD "/pixie"
