FROM golang:1.18.8-alpine3.16 as builder
WORKDIR /build

RUN apk add gcc git libc-dev

# Caching go dependencies.
COPY go.mod go.sum /
RUN go mod download
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN apk --no-cache add ca-certificates
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .

EXPOSE 8080
ENTRYPOINT CompileDaemon --build="go build -ldflags '-linkmode=external' cmd/main.go" --command=./main -log-prefix=false
