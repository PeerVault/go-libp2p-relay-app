FROM golang:1.12.9-alpine3.10

COPY relay.go go.* /app/
WORKDIR /app

RUN apk add --update git && \
    rm -rf /var/cache/apk/*

RUN go build -o /usr/bin/relay /app/relay.go

CMD relay
