# kubectl-curl

Kubectl plugin that executes a curl command from a dedicated Kubernetes pod.

Created for testing routing rules and network policies within a cluster, but it's also suitable for general API testing and data transfer.
The plugin simplifies these tasks by handling the deployment of a pod with the curl tool, then executing the curl command from within that pod.

## Installation

The latest version of kubectl-curl plugin can be installed using the following Go command:

```
go install github.com/michal-kopczynski/kubectl-curl@latest
```

Make sure that your `$GOPATH/bin` is included in your system's `PATH` to invoke the plugin from kubectl.

## Usage

The kubectl-curl command follows the standard syntax of curl. To execute plugin with default plugin options:
```
kubectl curl [curl options]
```

To execute plugin with custom plugin options:
```
kubectl curl [plugin options] -- [curl options]
```
The `--` ensures separation between kubectl-curl plugin options and the standard curl options.

## Examples
Execute a curl command using default plugin options:
```
kubectl curl -i http://httpbin/ip
```

Execute a curl command with custom plugin options:
```
kubectl curl -v -n foo -- -i http://httpbin/ip
```
