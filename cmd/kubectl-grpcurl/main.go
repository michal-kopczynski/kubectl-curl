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
		PluginKind:     plugin.Grpcurl,
		Version:        version,
		DefaultImage:   "fullstorydev/grpcurl:v1.8.9-alpine",
		DefaultPodName: "grpcurl",
		ExampleUsage: `# Execute a grpcurl command with default settings:
kubectl grpcurl -d {"greeting":"world"} -plaintext grpcbin:80 hello.HelloService.SayHello

# Execute a grpcurl command with custom plugin options. grpcurl options commes after '--':
kubectl grpcurl -v -n foo -- -d {"greeting":"world"} -plaintext grpcbin:80 hello.HelloService.SayHello`,
	}
	if err := cli.InitAndExecute(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
