# syntax = docker/dockerfile:experimental
FROM golang:1.21-bullseye as builder

RUN apt update && \
    apt -y install \
        apt-transport-https \
        ca-certificates

WORKDIR /usr/src/app
RUN groupadd --gid 2000 app && useradd -u 1000 -g app app

# used by the dev env
RUN go install github.com/cespare/reflex@latest

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/root/go/pkg go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -v -o webhook .

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=app:app --from=builder /usr/src/app/webhook /app/webhook
USER app
EXPOSE 8443
ENTRYPOINT ["/app/webhook"]
