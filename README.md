# Bitesize AuthZ Webhook

## Installation

`bitesize-authz-webhook` is packaged into docker container and can be found at
`geribatai/bitesize-authz-webhook`. We do not push a 'latest' container so reference the version number.
HTTPS termination is now supported.

## Configuration

- `/etc/kubernetes/authz-webhook/webhook.yaml`
```
clusters:
  - name: authz
    cluster:
      server: http://127.0.0.1:8888
current-context: webhook
contexts:
- context:
    cluster: authz
  name: webhook
```

- `/etc/kubernetes/authz-webhook/rules.hcl`
```
# see rules.hcl in examples
```

## Kubernetes configuration

### API server config

```
  --authorization-webhook-config-file=/etc/kubernetes/authz-webhook/webhook.yaml
  --authorization-mode=Webhook
```

### Run auth hook on the controller nodes (using a manifest)

- `/etc/kubernetes/manifests/kube-authz-webhook.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: kube-authz-webhook
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-authz-webhook
    image: geribatai/bitesize-authz-webhook:${VERSION}
    ports:
    - containerPort: 8888
      hostPort: 8888
    volumeMounts:
    - name: config
      mountPath: /etc/kubernetes/authz-webhook
      readOnly: true
    env:
    - name: LISTEN_PORT
      value: "8888"
    - name: RULES_CONFIG
      value: /etc/kubernetes/authz-webhook/rules.hcl
  volumes:
  - name: config
    hostPath:
      path: /etc/kubernetes/authz-webhook
```

### Environment variables

* `LISTEN_PORT` - Port webhook listens requests on (Default: 8080)
* `RULES_CONFIG` - path to `rules.hcl` file. (Default: rules.hcl in current
directory).
* `HTTPS_ENABLED` - (Default: false).
* `SSL_CERT` - (Default: cert.pem in current directory).
* `SSL_KEY` - (Default: key.pem in current directory).

### rules.hcl

Access rules are described in HCL format. Rules file is processed from the top,
and the first rule match found is returned as authorization status. If no match
is found, implicit deny rule is matched at the end.

```
access "allow" {
    user = "admin"
}

access "deny" {
    verb = "create"
}
```

Additional examples can be found in `examples` directory.


## Changelog
* 0.0.7 - Added debugging entries
* 0.0.6 - Added ACL support for groups
        - Added TLS support
* 0.0.5 - First open-source release. Supports HCL rules.
