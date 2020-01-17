.DEFAULT_GOAL := build
.PHONY: build
build: go build -o webhook .

.PHONY: start
build: webhook --tls-cert-file=/tls/tls.crt --tls-private-key-file=/tls/tls.key

.PHONY: test
test: go test ./...