package workerk8s

import (
	"context"
	"sync"

	"k8s.io/client-go/kubernetes"
)

// RunWorkers starts the specified number of worker goroutines to perform health checks on pods and collects their results.
func RunWorkers(ctx context.Context, clientset *kubernetes.Clientset, numWorkers int, namespace string) []string {
	results := make(chan string)
	var collectedResults []string
	var wg sync.WaitGroup

	// Start worker goroutines.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Worker(ctx, clientset, namespace, results)
		}()
	}

	// Collect results from the workers in a separate goroutine to avoid blocking.
	go func() {
		for result := range results {
			collectedResults = append(collectedResults, result)
		}
	}()

	// Wait for all workers to finish.
	wg.Wait()
	close(results) // Safe to close the channel here since all workers are done.

	return collectedResults
}
