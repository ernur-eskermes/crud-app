FROM golang:1.17-alpine AS builder

RUN go version

COPY . /github.com/ernur-eskermes/crud-app
WORKDIR /github.com/ernur-eskermes/crud-app

RUN go mod download
RUN GOOS=linux go build -o ./.bin/app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/ernur-eskermes/crud-app/.bin/app .

CMD ["./app"]