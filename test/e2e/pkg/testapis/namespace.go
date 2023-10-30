package testapis

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Namespace struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
	logger    *log.Logger
	name      string
}

func NewNamespace(clientset *kubernetes.Clientset, config *rest.Config, logger *log.Logger, name string) *Namespace {
	return &Namespace{
		clientset: clientset,
		config:    config,
		logger:    logger,
		name:      name,
	}
}

func (p *Namespace) Create() error {
	_, err := p.clientset.CoreV1().Namespaces().Create(context.TODO(), &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.name,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}

	p.logger.Printf("Namespace \"%s\" created successfully.\n", p.name)
	return nil
}

func (p *Namespace) Delete() error {
	err := p.clientset.CoreV1().Namespaces().Delete(context.TODO(), p.name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete namespace: %v", err)
	}

	return nil
}
