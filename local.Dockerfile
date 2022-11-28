FROM golang:1.18.8-alpine3.16 as builder
WORKDIR /build

RUN apk add gcc git libc-dev

# Caching go dependencies.
COPY go.mod go.sum /
RUN go mod download
RUN go install -mod=mod github.com/githubnemo/CompileDaemon

COPY . .

EXPOSE 8080
ENTRYPOINT CompileDaemon --build="go build cmd/main.go" --command=./main
