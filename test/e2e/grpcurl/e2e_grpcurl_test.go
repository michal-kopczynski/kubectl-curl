package e2e

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/michal-kopczynski/kubectl-curl/pkg/apis"
	"github.com/michal-kopczynski/kubectl-curl/pkg/plugin"
	"github.com/michal-kopczynski/kubectl-curl/test/e2e/pkg/testapis"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	testNamespaceName = "kubectl-grpcurl-test"
	httpbinPodName    = "grpcbin"
	httpbinImage      = "kong/grpcbin:latest"
)

type TestState struct {
	clientset      *kubernetes.Clientset
	testNamespace  *testapis.Namespace
	httpbinPod     *apis.Pod
	httpbinService *testapis.Service
}

// Requires Kubernetes cluster which can be created using for example Minikube

func setup(t *testing.T) *TestState {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.SetOutput(io.Discard)

	kubeconfig := plugin.GetKubeconfig("")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating clientset: %v", err)
	}

	testNamespace := testapis.NewNamespace(
		clientset,
		config,
		logger,
		testNamespaceName)
	httpbinPod := apis.NewPod(
		clientset,
		config,
		logger,
		httpbinImage,
		testNamespaceName,
		httpbinPodName,
		[]string{},
		9000)
	httpbinService := testapis.NewService(
		clientset,
		config,
		logger,
		testNamespaceName,
		httpbinPodName,
		80,
		9000)

	if err := testNamespace.Create(); err != nil {
		t.Fatalf("Error creating test namespace: %v", err)
	}

	if err := httpbinService.Create(); err != nil {
		t.Fatalf("Error creating httpbin service: %v", err)
	}

	if err := httpbinPod.Create(); err != nil {
		t.Fatalf("Error creating httpbin pod: %v", err)
	}

	if err := httpbinPod.WaitForReady(30 * time.Second); err != nil {
		t.Fatalf("Error waiting for httpbin pod readiness: %v", err)
	}

	return &TestState{
		clientset:      clientset,
		testNamespace:  testNamespace,
		httpbinPod:     httpbinPod,
		httpbinService: httpbinService,
	}
}

func teardown(t *testing.T, testState *TestState) {
	if err := testState.httpbinPod.Delete(); err != nil {
		t.Fatalf("Error deleting httpbin pod: %v", err)
	}

	if err := testState.httpbinService.Delete(); err != nil {
		t.Fatalf("Error deleting httpbin service: %v", err)
	}

	if err := testState.testNamespace.Delete(); err != nil {
		t.Fatalf("Error deleting test namespace: %v", err)
	}
}

func TestKubectlCurl(t *testing.T) {
	testState := setup(t)
	defer teardown(t, testState)

	tests := []struct {
		name             string
		curlArgs         string
		expectToPass     bool
		expectedInOutput []string
	}{
		{
			name:             "Test default plugin options - TLS failure",
			curlArgs:         `-d {"greeting":"world"} grpcbin.` + testNamespaceName + `.svc.cluster.local:80 hello.HelloService.SayHello`,
			expectToPass:     false,
			expectedInOutput: []string{"Failed to dial target host", "first record does not look like a TLS handshak"},
		},
		{
			name:             "Test default plugin and grpcurl options",
			curlArgs:         `-d {"greeting":"world"} -plaintext grpcbin.` + testNamespaceName + `.svc.cluster.local:80 hello.HelloService.SayHello`,
			expectToPass:     true,
			expectedInOutput: []string{"hello world"},
		},
		{
			name:             "Test custom namespace option in plugin options",
			curlArgs:         `-n ` + testNamespaceName + ` -- -d {"greeting":"world"} -plaintext grpcbin:80 hello.HelloService.SayHello`,
			expectToPass:     true,
			expectedInOutput: []string{"hello world"},
		},
		{
			name:             "Test verbose option in plugin options",
			curlArgs:         `-v -n ` + testNamespaceName + ` -- -d {"greeting":"world"} -plaintext grpcbin:80 hello.HelloService.SayHello`,
			expectToPass:     true,
			expectedInOutput: []string{"Using kubeconfig", "Executing: grpcurl -d {\"greeting\":\"world\"} -plaintext grpcbin:80 hello.HelloService.SayHello", "hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandArgs := append([]string{"grpcurl"}, strings.Split(tt.curlArgs, " ")...)
			cmd := exec.Command("kubectl", commandArgs...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if tt.expectToPass {
				if err != nil {
					t.Fatalf("Failed to run kubectl grpcurl: %v", err)
				}
			}

			output := out.String()
			t.Log(output)

			for _, expected := range tt.expectedInOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q. Output:\n%s", expected, output)
				}
			}
		})
	}
}
