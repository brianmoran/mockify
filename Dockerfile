FROM golang:1
COPY app /app
WORKDIR app
RUN ls -alh
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/mockify.go
WORKDIR cmd
COPY main /main
CMD ["main"]
