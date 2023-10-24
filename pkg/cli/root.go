package cli

import (
	"io"
	"log"
	"os"

	"github.com/michal-kopczynski/kubectl-curl/pkg/plugin"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func RootCmd(version string) *cobra.Command {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	opts := &plugin.Opts{
		Kubeconfig: "",
		Namespace:  "default",
		PodName:    "curl",
		Cleanup:    false,
		Verbose:    false,
		Timeout:    30,
	}

	cmd := &cobra.Command{
		Use: `kubectl curl [curl options]
  kubectl curl [plugin options] -- [curl options]`,
		Short: "Executes a curl command from a dedicated Kubernetes pod",
		Example: `# Execute a curl command using default plugin options.
kubectl curl -i http://httpbin/ip

# Execute a curl command with custom plugin options. curl options commes after '--'.
kubectl curl -v -n foo -- -i http://httpbin/ip`,
		SilenceUsage: true,
		Version:      version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Verbose {
				logger.SetOutput(io.Discard)
			}
			return plugin.RunPlugin(logger, opts, args)
		},
	}

	if slices.Contains(os.Args, "--") || slices.Contains(os.Args, "--help") {
		cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", opts.Kubeconfig, "path to kubeconfig file")
		cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", opts.Namespace, "namespace in which curl pod will be created")
		cmd.Flags().StringVar(&opts.PodName, "name", opts.PodName, "curl pod name")
		cmd.Flags().BoolVarP(&opts.Cleanup, "cleanup", "c", opts.Cleanup, "delete curl pod at the end")
		cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", opts.Verbose, "explain what is being done")
		cmd.Flags().IntVarP(&opts.Timeout, "timeout", "t", opts.Timeout, "the timeout of plugin operations in seconds")
	} else {
		cmd.DisableFlagParsing = true
	}

	return cmd
}

func InitAndExecute(version string) error {
	if err := RootCmd(version).Execute(); err != nil {
		return err
	}
	return nil
}
