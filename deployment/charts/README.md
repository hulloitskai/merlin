# charts

_[Helm](https://helm.sh) charts for [`merlin`][merlin]._

## Installation

You can install these charts from the repository located at `https://charts.stevenxie.me`.

```bash
## Add the repository.
helm repo add stevenxie https://charts.stevenxie.me

## Install the chart.
helm install -f values.yaml -n merlin stevenxie/merlin
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

[merlin]: https://github.com/stevenxie/merlin
