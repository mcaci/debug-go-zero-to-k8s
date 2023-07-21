FROM golang:1.10.0-alpine

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . /source

ENTRYPOINT ["sh"]
# go/bin/dlv attach $PID