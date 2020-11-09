FROM golang:1-alpine as builder

RUN apk update && \
    apk add git && \
    rm -rf /var/cache/apk/*

RUN go get github.com/gorilla/mux \
    github.com/json-iterator/go \
    gopkg.in/yaml.v2

WORKDIR /go/src/mockify/

COPY . /go/src/mockify/

RUN CGO_ENABLED=0 go build -v -o mockify ./app/cmd/mockify.go

FROM alpine:latest

EXPOSE 8001

RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /root/

COPY config/* config/
COPY --from=builder /go/src/mockify/mockify .

CMD ./mockify
