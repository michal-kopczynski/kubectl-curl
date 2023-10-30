package plugin

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/michal-kopczynski/kubectl-curl/pkg/apis"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type PluginKind string

const (
	Curl    PluginKind = "curl"
	Grpcurl            = "grpcurl"
)

func (p PluginKind) String() string {
	return string(p)
}

type Opts struct {
	Kubeconfig string
	Image      string
	Namespace  string
	PodName    string
	Cleanup    bool
	Verbose    bool
	Timeout    int
}

func GetKubeconfig(kubeconfig string) string {
	if kubeconfig != "" {
		return kubeconfig
	} else if kubeconfigEnv, exists := os.LookupEnv("KUBECONFIG"); exists {
		return kubeconfigEnv
	} else {
		if home := homedir.HomeDir(); home != "" {
			return filepath.Join(home, ".kube", "config")
		}
	}
	return ""
}

func RunPlugin(kind PluginKind, logger *log.Logger, opts *Opts, args []string) error {
	curlCommand := kind.String() + " " + strings.Join(args, " ")
	timeout := time.Duration(opts.Timeout) * time.Second

	kubeconfig := GetKubeconfig(opts.Kubeconfig)
	logger.Printf("Using kubeconfig: %s\n", kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating clientset: %w", err)
	}

	pod := apis.NewPod(
		clientset,
		config,
		logger,
		opts.Image,
		opts.Namespace,
		opts.PodName,
		[]string{"sleep", "infinity"},
		0)

	podExists, err := pod.IsCreated()
	if err != nil {
		return fmt.Errorf("error checking if \"%s\" exists: %w", opts.PodName, err)
	}

	if !podExists {
		if err := pod.Create(); err != nil {
			return fmt.Errorf("error creating \"%s\" pod: %w", opts.PodName, err)
		}
	} else {
		logger.Printf("Pod \"%s\" already exists.", opts.PodName)
	}

	if err := pod.WaitForReady(timeout); err != nil {
		return fmt.Errorf("error waiting for \"%s\" readiness: %w", opts.PodName, err)
	}

	output, err := pod.ExecuteCommand(curlCommand, timeout)
	if err != nil {
		return fmt.Errorf("error executing command inside \"%s\" pod: %w", opts.PodName, err)
	}
	fmt.Println(output)

	if opts.Cleanup {
		if err := pod.Delete(); err != nil {
			return fmt.Errorf("error deleting \"%s\" pod: %w", opts.PodName, err)
		}
	}

	return nil
}
