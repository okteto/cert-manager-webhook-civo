# Civo ACME DNS Webhook

Cert Manager Civo Webhook performing ACME challenge using DNS record

## Installation

```
$ helm install civo chart/certmanager-civo --namespace cert-manager
```

# How to Use

## Secret
```
kubectl create secret generic=
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
              name: okteto-dns
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