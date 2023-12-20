package workerk8s

import (
	"k8s.io/client-go/kubernetes"
)

// RunWorkers starts the specified number of worker goroutines to perform health checks on pods and collects their results.
func RunWorkers(clientset *kubernetes.Clientset, numWorkers int, namespace string) []string {
	results := make(chan string)
	var collectedResults []string

	// Start worker goroutines.
	for i := 0; i < numWorkers; i++ {
		go Worker(clientset, namespace, results)
	}

	// Collect results from the workers.
	for i := 0; i < numWorkers; i++ {
		collectedResults = append(collectedResults, <-results)
	}

	// It's important to close the channel to avoid a goroutine leak.
	close(results)

	return collectedResults
}
