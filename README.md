# kubectl-curl and kubectl-grpcurl

Kubectl plugins that execute curl and grpcurl commands from a dedicated Kubernetes pod.

Created for testing routing rules and network policies within a cluster, but they are also suitable for general API testing and data transfer.
These plugins simplify these tasks by handling the deployment of pods with the curl/grpcurl tools and then executing the curl/grpcurl commands from within those pods.

## Installation

### Using release binaries

The latest release binaries for kubectl-curl and kubectl-grpcurl plugins can be downloaded from the [GitHub releases page](https://github.com/michal-kopczynski/kubectl-curl/releases/latest).

### Using go install

The latest versions of kubectl-curl and kubectl-grpcurl plugins can be installed using the following Go command:

```
go install github.com/michal-kopczynski/kubectl-curl/...@latest
```

Make sure that your `$GOPATH/bin` is included in your system's `PATH` to invoke the plugins from kubectl.

To install only kubectl curl plugin:
```
go install github.com/michal-kopczynski/kubectl-curl/cmd/kubectl-curl@latest
```

To install only kubectl grpcurl plugin:
```
go install github.com/michal-kopczynski/kubectl-curl/cmd/kubectl-grpcurl@latest
```

## Usage

The kubectl-curl and kubectl-grpcurl commands follow the standard syntax of curl/grpcurl.

To execute plugins with default options (deploy curl/grpcurl pod in default namespace):
```
kubectl curl [curl options]
kubectl grpcurl [grpcurl options]
```

To execute plugins with custom options:
```
kubectl curl [plugin flags] -- [curl options]
kubectl grpcurl [plugin flags] -- [grpcurl options]
```
The `--` ensures separation between kubectl-curl/kubectl-grpcurl plugins flags and the standard curl/grpcurl options.

## Examples
Execute a curl/grpcurl command using default options:
```
kubectl curl -i http://httpbin/ip
kubectl grpcurl -d '{"greeting":"world"}' -plaintext grpcbin:80 hello.HelloService.SayHello
```

Execute a curl/grpcurl command with custom options i.e:
```
kubectl curl --verbose --namespace foo -- -i http://httpbin/ip
kubectl grpcurl --verbose --namespace foo -- -d '{"greeting":"world"}' -plaintext grpcbin:80 hello.HelloService.SayHello
```
