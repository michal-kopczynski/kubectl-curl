package main

import (
	"fmt"
	"os"

	"github.com/michal-kopczynski/kubectl-curl/pkg/cli"
	"github.com/michal-kopczynski/kubectl-curl/pkg/plugin"
)

var version string

func main() {
	c := cli.Config{
		PluginKind:     plugin.Curl,
		Version:        version,
		DefaultImage:   "curlimages/curl:8.4.0",
		DefaultPodName: "curl",
		ExampleUsage: `# Execute a curl command with default settings.
kubectl curl -i http://httpbin/ip

# Execute a curl command with custom plugin options. curl options commes after '--'.
kubectl curl -v -n foo -- -i http://httpbin/ip`,
	}
	if err := cli.InitAndExecute(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
