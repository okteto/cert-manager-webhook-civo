# Civo ACME DNS Webhook

An ACME Webhook solver for [CIVO](https://civo.com)

## Installation

```
$ helm install civo chart/civo-acme --namespace cert-manager
```

# How to Use

## Secret
```
kubectl create secret generic dns --from-literal=key=<YOUR_CIVO_TOKEN>
```

# Issuer
```
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: civoissuer
spec:
  acme:
    email: example@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    server: https://acme-v02.api.letsencrypt.org/directory
    solvers:
    - dns01:
        webhook:
          solverName: "civo"
          groupName: civo.webhook.okteto.com
          config:
            apiKeySecretRef:
              key: key
              name: dns
```

## Certificate
```
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: my-certificate
spec:
  dnsNames:
  - '*.example.com'
  issuerRef:
    kind: Issuer
    name: civo
  secretName: wildcard-example-com-tls
```

# Contributing
If you want to get involved, we'd love to receive a pull request, issues, or an offer to help over at the [#KUBE100](https://app.slack.com/client/TKW8H5MBK/CMVCKMCN5) channel in the Civo-Community slack or at the [#Okteto](https://kubernetes.slack.com/messages/CM1QMQGS0/) channel in the Kubernetes slack.

Maintainers:
- [Ramiro Berrelleza](https://twitter.com/rberrelleza)
- [Pablo Chico de Guzman](https://twitter.com/pchico83)

Please see the [contribution guidelines](CONTRIBUTING.md)
