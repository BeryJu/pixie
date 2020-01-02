FROM golang:latest AS builder
WORKDIR $GOPATH/src/github.com/BeryJu/pixie
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go get github.com/gobuffalo/packr/v2/packr2 && \
    packr2 && \
    go build -v -o /go/bin/pixie

FROM scratch
COPY --from=builder /go/bin/pixie /pixie
EXPOSE 8080
WORKDIR /web-root
CMD [ "/pixie" ]
ENTRYPOINT [ "/pixie" ]
