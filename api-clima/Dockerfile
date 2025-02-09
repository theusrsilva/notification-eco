FROM golang:1.23.5-alpine3.20

RUN apk add --no-cache bash
ENV PATH="$PATH:/bin/bash"

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]