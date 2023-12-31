package workerk8s

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// RunWorkers starts the specified number of worker goroutines to perform health checks on pods and collects their results.
// It returns a channel to receive the results and a function to trigger a graceful shutdown.
func RunWorkers(ctx context.Context, clientset *kubernetes.Clientset, namespace string, workerCount int) (<-chan string, func()) {
	results := make(chan string)
	var wg sync.WaitGroup

	shutdownCtx, cancelFunc := context.WithCancel(ctx)

	// Start the specified number of worker goroutines.
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			// We now use the package-level Logger, enhanced with additional fields.
			SetLogger(Logger.With(zap.Int("workerIndex", workerIndex)))
			Worker(shutdownCtx, clientset, namespace, results)
		}(i)
	}

	// Shutdown function to be called to initiate a graceful shutdown.
	shutdown := func() {
		// Signal all workers to stop by cancelling the context.
		cancelFunc()

		// Wait for all workers to finish.
		go func() {
			wg.Wait()
			close(results)
		}()
	}

	return results, shutdown
}
