# syntax = docker/dockerfile:experimental
FROM golang:1.18-buster as builder

RUN apt update && \
    apt -y install \
        apt-transport-https \
        ca-certificates

WORKDIR /usr/src/app

# used by the dev env
RUN go get -u github.com/cespare/reflex

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/root/go/pkg go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -v -o webhook .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/webhook /app/webhook
EXPOSE 443
ENTRYPOINT ["/app/webhook"]