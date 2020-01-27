module github.com/okteto/civo-acme

go 1.13

require (
	github.com/civo/civogo v0.0.0-20200123135111-b3aba767c3d7
	github.com/jetstack/cert-manager v0.12.0
	github.com/okteto/civogo v0.0.0-20200116195624-aa4f756bebb9
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.4

replace github.com/civo/civogo => github.com/rberrelleza/civogo v0.0.0-20200126190228-5c6fe06b4d2b
