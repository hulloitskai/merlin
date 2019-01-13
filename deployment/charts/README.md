# charts

This directory contains [Helm](https://helm.sh) charts for this project.

You can install these from the packaged `.tar.gz` files
[from the latest Github release](https://github.com/stevenxie/merlin/releases).

For example, to install `merlin@0.4.3`, you can simply run:

```bash
helm install -n merlin -f values.yaml \
  https://github.com/stevenxie/merlin/releases/download/v0.4.3/merlin-0.1.0.tgz
```

## Configuration

See
[`merlin/values.yaml`](https://github.com/stevenxie/merlin/blob/master/deployment/charts/merlin/values.yaml)
for an the default `values.yaml` configuration.

To install `merlin` for production, one should have an Ingress controller in
the target namespace, and configure a `values.yaml` with an appropriate
`ingress.host` value:

```yaml
ingress:
  host: merlin.stevenxie.me # example
```
