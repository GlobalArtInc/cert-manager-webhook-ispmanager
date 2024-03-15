# cert-manager-webhook-ispmanager

Cert-manager ACME DNS webhook provider for ISPManager.

## Installing

To install with helm, run:

```bash
$ helm repo add globalart https://globalartinc.github.io/helm-charts
$ helm upgrade -n cert-manager --install cert-manager-webhook-ispmanager globalart/cert-manager-webhook-ispmanager
```


### Issuer/ClusterIssuer

An example issuer:

```yaml
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: letsencrypt
  namespace: default
spec:
  acme:
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    email: webmaster@globalart.dev
    solvers:
    - dns01:
        webhook:
          groupName: acme.ispmanager.com
          solverName: ispmanager-provider
          config:
            panelUrl: "your_panel_url"
            user: "username"
            password: "password"
```

And then you can issue a cert:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: sel-letsencrypt-crt
  namespace: default
spec:
  secretName: example-com-tls
  commonName: example.com
  issuerRef:
    name: letsencrypt-staging
    kind: Issuer
  dnsNames:
  - example.com
  - www.example.com
```

## Development

### Running the test suite

You can run the test suite with:

1. Fill in the appropriate values in `testdata/ispmanager/config.json` 
2. Change `dns.SetDNSServer("127.0.0.1:53")` on your DNS Server

```bash
$ TEST_ZONE_NAME=example.com. make test
```