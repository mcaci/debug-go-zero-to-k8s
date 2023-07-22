FROM golang:1.20.0-alpine

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENV APP_PATH=/home/mcaci/go/src/github.com/mcaci/debug-go-zero-to-k8s/
COPY . $APP_PATH

ENTRYPOINT ["sh"]