# syntax = docker/dockerfile:experimental
FROM okteto/golang:1.17 as builder

WORKDIR /usr/src/app

# used by the dev env
RUN go get -u github.com/cespare/reflex

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/root/go/pkg go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -v -o webhook .

FROM debian:buster
RUN apt-get update \
        && apt-get install -y ca-certificates

COPY --from=builder /usr/src/app/webhook /app/webhook
EXPOSE 443
ENTRYPOINT ["/app/webhook"]