#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./app/cmd/mockify.go
docker rmi -f mockify
docker build -t mockify -f Dockerfile.scratch .