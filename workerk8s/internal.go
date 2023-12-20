package workerk8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewKubernetesClient creates a new Kubernetes client using the in-cluster configuration.
// This is typically used when the application itself is running within a Kubernetes cluster.
func NewKubernetesClient() (*kubernetes.Clientset, error) {
	// Get the in-cluster config.
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Create a clientset based on the in-cluster config.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
