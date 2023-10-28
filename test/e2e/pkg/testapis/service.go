package testapis

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Service struct {
	clientset  *kubernetes.Clientset
	config     *rest.Config
	logger     *log.Logger
	namespace  string
	name       string
	port       int32
	targetPort int32
}

func NewService(clientset *kubernetes.Clientset, config *rest.Config, logger *log.Logger, namespace string, name string, port int32, targetPort int32) *Service {
	return &Service{
		clientset:  clientset,
		config:     config,
		logger:     logger,
		namespace:  namespace,
		name:       name,
		port:       port,
		targetPort: targetPort,
	}
}

func (p *Service) Create() error {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.name,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": p.name,
			},
			Ports: []apiv1.ServicePort{
				{
					Port:       p.port,
					TargetPort: intstr.FromInt32(p.targetPort),
				},
			},
		},
	}

	_, err := p.clientset.CoreV1().Services(p.namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	p.logger.Printf("Service \"%s\" created successfully in namespace \"%s\".\n", p.name, p.namespace)
	return nil
}

func (p *Service) Delete() error {
	err := p.clientset.CoreV1().Services(p.namespace).Delete(context.TODO(), p.name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil

}
