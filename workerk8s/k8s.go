package workerk8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewKubernetesClient creates and returns a new Kubernetes client.
func NewKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// Worker interacts with the Kubernetes API and sends results to a channel.
func Worker(clientset *kubernetes.Clientset, results chan<- string) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		results <- fmt.Sprintf("Error: %v", err)
		return
	}
	results <- fmt.Sprintf("Worker processed %d pods", len(pods.Items))
}

// RunWorkers starts a number of worker goroutines and collects their results.
func RunWorkers(clientset *kubernetes.Clientset, numWorkers int) []string {
	results := make(chan string)
	var collectedResults []string

	for i := 0; i < numWorkers; i++ {
		go Worker(clientset, results)
	}

	for i := 0; i < numWorkers; i++ {
		collectedResults = append(collectedResults, <-results)
	}

	return collectedResults
}
