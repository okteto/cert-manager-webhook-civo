# syntax = docker/dockerfile:experimental
FROM golang:1.13 as builder_deps

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM builder_deps as builder
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -v -o webhook .

FROM alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

COPY --from=builder /app/webhook /app/webhook
EXPOSE 443
ENTRYPOINT ["/app/webhook"]