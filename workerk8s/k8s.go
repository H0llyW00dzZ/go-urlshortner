package workerk8s

import (
	"k8s.io/client-go/kubernetes"
)

// RunWorkers starts a number of worker goroutines and collects their results.
func RunWorkers(clientset *kubernetes.Clientset, numWorkers int, namespace string) []string {
	results := make(chan string)
	var collectedResults []string

	for i := 0; i < numWorkers; i++ {
		go Worker(clientset, namespace, results)
	}

	for i := 0; i < numWorkers; i++ {
		collectedResults = append(collectedResults, <-results)
	}

	return collectedResults
}
