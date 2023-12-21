package workerk8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewKubernetesClient creates a new Kubernetes client using the in-cluster configuration.
// This is typically used when the application itself is running within a Kubernetes cluster.
func NewKubernetesClient() (*kubernetes.Clientset, error) {
	// Use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		// If running outside the cluster, use the kubeconfig file.
		kubeconfig := filepath.Join(os.Getenv(HOME), kube, Config)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf(errconfig, err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf(cannotcreatek8s, err)
	}

	return clientset, nil
}
