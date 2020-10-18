.PHONY: build
build: 
	go build -o webhook

.PHONY: start
start: 
	webhook --tls-cert-file=/tls/tls.crt --tls-private-key-file=/tls/tls.key

.PHONY: test
test: 
	go test ./...