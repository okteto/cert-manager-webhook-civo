# CIVO DNS webhook for cert-manager

A webhook to use [CIVO DNS](https://civo.com) as a DNS issuer for [cert-manager](https://github.com/jetstack/cert-manager).

## Installation

```
$ helm install webhook-civo https://storage.googleapis.com/charts.okteto.com/cert-manager-webhook-civo-0.1.0.tgz --namespace=cert-manager
```

# How to Use

## Secret
```
kubectl create secret generic dns --namespace=cert-manager --from-literal=key=<YOUR_CIVO_TOKEN>
```

# Issuer

> Update email to match yours 
```
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: civo
spec:
  acme:
    email: YOUR_EMAIL
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

> Update the dnsNames to match yours 

```
apiVersion: cert-manager.io/v1
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
