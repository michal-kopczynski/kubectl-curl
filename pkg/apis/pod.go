package apis

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

type Pod struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
	logger    *log.Logger
	image     string
	namespace string
	name      string
	command   []string
	port      int32
}

func NewPod(clientset *kubernetes.Clientset, config *rest.Config, logger *log.Logger, image string, namespace string, name string, command []string, port int32) *Pod {
	return &Pod{
		clientset: clientset,
		config:    config,
		logger:    logger,
		image:     image,
		namespace: namespace,
		name:      name,
		command:   command,
		port:      port,
	}
}

func (p *Pod) IsCreated() (bool, error) {
	podsClient := p.clientset.CoreV1().Pods(p.namespace)

	_, err := podsClient.Get(context.TODO(), p.name, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *Pod) Create() error {
	podsClient := p.clientset.CoreV1().Pods(p.namespace)

	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.name,
			Labels: map[string]string{
				"app": p.name,
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  p.name,
					Image: p.image,
				},
			},
		},
	}

	if len(p.command) != 0 {
		pod.Spec.Containers[0].Command = p.command
	}

	if p.port != 0 {
		pod.Spec.Containers[0].Ports = []apiv1.ContainerPort{
			{
				ContainerPort: p.port,
			},
		}
	}

	_, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {

		return fmt.Errorf("failed to create pod: %w", err)
	}

	p.logger.Printf("Pod \"%s\" created successfully in namespace \"%s\".\n", p.name, p.namespace)

	return nil
}

func (p *Pod) WaitForReady(timeout time.Duration) error {
	watch, err := p.clientset.CoreV1().Pods(p.namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + p.name,
	})
	if err != nil {
		return fmt.Errorf("Failed to watch pod: %w", err)
	}

	timeoutChan := time.After(timeout)
	p.logger.Println("Waiting for pod to be ready...")

	for {
		select {
		case event := <-watch.ResultChan():
			pod, ok := event.Object.(*apiv1.Pod)
			if !ok {
				return fmt.Errorf("unexpected type in watch event")
			}
			if pod.Status.Phase == apiv1.PodRunning {
				p.logger.Println("Pod is now running.")
				return nil
			}
		case <-timeoutChan:
			return fmt.Errorf("timed out waiting for pod to be ready")
		}
	}
}

func (p *Pod) ExecuteCommand(command string, timeout time.Duration) (string, error) {
	execRequest := p.clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(p.name).
		Namespace(p.namespace).
		SubResource("exec").
		VersionedParams(&apiv1.PodExecOptions{
			Container: p.name,
			Command:   strings.Split(command, " "),
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.config, "POST", execRequest.URL())
	if err != nil {
		return "", fmt.Errorf("Failed to initialize command executor: %w", err)
	}

	output := &strings.Builder{}
	errorOutput := &strings.Builder{}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p.logger.Printf("Executing: %s", command)

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: output,
		Stderr: errorOutput,
	})
	if err != nil {
		p.logger.Printf("Command failed: %s\n", err)
		return errorOutput.String(), nil
	}

	p.logger.Println("Command executed successfully. Output:")

	return output.String(), nil
}

func (p *Pod) Delete() error {
	podsClient := p.clientset.CoreV1().Pods(p.namespace)

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	err := podsClient.Delete(context.TODO(), p.name, *deleteOptions)
	if err != nil {
		return fmt.Errorf("Failed to delete pod: %w", err)
	}

	p.logger.Println("Pod deleted successfully.")
	return nil
}
