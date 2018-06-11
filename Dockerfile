FROM golang:1 as builder

RUN go get github.com/gorilla/mux
RUN go get github.com/sirupsen/logrus

COPY . /go/src/mockify/
WORKDIR /go/src/mockify/

RUN CGO_ENABLED=0 go build -v -o mockify app/cmd/mockify.go

FROM alpine:latest

EXPOSE 8001

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root/

COPY config/* config/
COPY --from=builder /go/src/mockify/mockify .

CMD ./mockify
