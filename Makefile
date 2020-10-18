.PHONY: build
build: 
	go build -o webhook

.PHONY: start
start: build
	./webhook --tls-cert-file=/tls/tls.crt --tls-private-key-file=/tls/tls.key


watch:
	ulimit -n 1000
	reflex -s -r '\.go$$' make start

.PHONY: test
test: 
	TEST_ZONE_NAME=example.com. go test ./...