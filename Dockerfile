# syntax = docker/dockerfile:experimental
FROM golang:1.15 as builder_deps

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM builder_deps as builder
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -v -o webhook .

FROM debian:buster
RUN apt-get update \
        && apt-get install -y ca-certificates

COPY --from=builder /app/webhook /app/webhook
EXPOSE 443
ENTRYPOINT ["/app/webhook"]