FROM golang:1.23.5-alpine3.20

RUN apk add --no-cache bash
ENV PATH="$PATH:/bin/bash"

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src
RUN apk add --no-cache gcc musl-dev sqlite-dev

ENTRYPOINT ["tail", "-f", "/dev/null"]