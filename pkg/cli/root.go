package cli

import (
	"io"
	"log"
	"os"
	"slices"

	"github.com/michal-kopczynski/kubectl-curl/pkg/plugin"
	"github.com/spf13/cobra"
)

type Config struct {
	PluginKind     plugin.PluginKind
	Version        string
	DefaultImage   string
	DefaultPodName string
	ExampleUsage   string
}

func RootCmd(config Config) *cobra.Command {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	opts := &plugin.Opts{
		Kubeconfig: "",
		Image:      config.DefaultImage,
		Namespace:  "default",
		PodName:    config.DefaultPodName,
		Cleanup:    false,
		Verbose:    false,
		Timeout:    30,
	}

	pluginName := config.PluginKind.String()
	cmd := &cobra.Command{
		Use: `kubectl ` + pluginName + ` [` + pluginName + ` options]
  kubectl ` + pluginName + ` [plugin flags] -- [` + pluginName + ` options]`,
		Short:        "Executes a " + pluginName + " command from a dedicated Kubernetes pod",
		Example:      config.ExampleUsage,
		SilenceUsage: true,
		Version:      "kubect-" + config.PluginKind.String() + " version: " + config.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Verbose {
				logger.SetOutput(io.Discard)
			}
			return plugin.RunPlugin(config.PluginKind, logger, opts, args)
		},
	}

	cmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)

	cmd.DisableFlagsInUseLine = true

	if slices.Contains(os.Args, "--") || slices.Contains(os.Args, "--help") || slices.Contains(os.Args, "--version") {
		cmd.Flags().StringVar(&opts.Kubeconfig, "kubeconfig", opts.Kubeconfig, "path to kubeconfig file")
		cmd.Flags().StringVarP(&opts.Image, "image", "i", opts.Image, "docker image with "+pluginName+" tool")
		cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", opts.Namespace, "namespace in which "+pluginName+" pod will be created")
		cmd.Flags().StringVar(&opts.PodName, "name", opts.PodName, pluginName+" pod name")
		cmd.Flags().BoolVarP(&opts.Cleanup, "cleanup", "c", opts.Cleanup, "delete "+pluginName+" pod at the end")
		cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", opts.Verbose, "explain what is being done")
		cmd.Flags().IntVarP(&opts.Timeout, "timeout", "t", opts.Timeout, "the timeout of plugin operations in seconds")
	} else {
		cmd.DisableFlagParsing = true
	}

	return cmd
}

func InitAndExecute(config Config) error {
	if err := RootCmd(config).Execute(); err != nil {
		return err
	}
	return nil
}
