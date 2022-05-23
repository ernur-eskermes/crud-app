FROM alpine:latest

RUN apk update & apk --no-cache add ca-certificates
WORKDIR /root/

CMD ["./app"]