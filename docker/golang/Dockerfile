FROM golang:1.14-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh libc-dev gcc openssl

WORKDIR /app