package main

import (
	"fmt"
	"os"

	"github.com/michal-kopczynski/kubectl-curl/pkg/cli"
)

var version string

func main() {
	if err := cli.InitAndExecute(version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
