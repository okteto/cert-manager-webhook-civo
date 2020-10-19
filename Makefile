watch:
	ulimit -n 1000
	reflex -s -r '\.go$$' make start

.PHONY: start
start: build
	./webhook --tls-cert-file=/tls/tls.crt --tls-private-key-file=/tls/tls.key

.PHONY: build
build: 
	go build -o webhook

.PHONY: test
test: 
	TEST_ZONE_NAME=example.com. go test ./...

.PHONY: readme
readme:
	helm-docs

.PHONY: chart
chart:
	okteto build -t okteto/civo-webhook:${TAG}
	helm package chart/cert-manager-webhook-civo --app-version ${TAG}
	mkdir -p output
	mv *.tgz output/

release:
	curl https://charts.okteto.com/index.yaml > old-index.yaml
	helm repo index --merge old-index.yaml --url https://charts.okteto.com output
	gsutil cp output/index.yaml gs://charts.okteto.com/
	gsutil cp output/*.tgz gs://charts.okteto.com/