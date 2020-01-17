# Contributing

If you submit a pull request, please keep the following guidelines in mind:

1. Code should be `go fmt` compliant.
2. Types, structs and funcs should be documented.
3. Tests pass.

## Getting set up

Since this webhook depends on cert-manager, certificates and other Kubernetes resources, it's a lot easier to develop it directly in Kubernetes with `okteto`. 

Clone the repo:
```sh
git clone https://github.com/okteto/civo-acme
```

Deploy the latest version of webhook via Helm 3 in your cluster:
`helm install dev charts/civo-acme`

Deploy your development environment with [`okteto`](https://github.com/okteto/okteto):
    `okteto up`

Build the latest version:
`okteto> make build`

Start the webhook process:
`okteto> make start`

Okteto will automatically synchronize your changes, so rebuilding and restarting the process is all you need to do when validating a change.

## Running tests

When working on code in this repository, tests can be run via:

```sh
make test
```

## Validate your changes end to end
Easiest way is to create and delete a certificate from the command line. The [samples](samples) directory has a certificate and an issuer to help you get started.
