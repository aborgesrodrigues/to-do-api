FROM golang:1.24-alpine as builder
WORKDIR /build
RUN apk add gcc git libc-dev

# Caching go dependencies.
COPY go.mod go.sum /
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-linkmode=external' -o server cmd/main.go

# Must match the alpine version used in the builder stage for Dynatrace
# to instrument correctly.
FROM alpine:latest
WORKDIR /usr/local/bin/
EXPOSE 8080
RUN apk add --no-cache ca-certificates libc6-compat

COPY --from=builder /build/server .
# Prefer ENTRYPOINT over CMD and exec form over shell form for reasons documented at:
# https://docs.docker.com/engine/reference/builder/#entrypoint
ENTRYPOINT ["./server"]
