// Package workerk8s provides a simplified interface to interact with the Kubernetes API.
// It encapsulates the creation of a Kubernetes client and provides functionality
// to execute concurrent operations using worker goroutines.
//
// The package is designed to be used in environments where Kubernetes is the orchestrator
// and the application is running as a pod within the Kubernetes cluster. It leverages
// the in-cluster configuration to create a client that can interact with the API server.
//
// # Functions
//
//   - NewKubernetesClient: Establishes a new connection to the Kubernetes API server using
//     in-cluster configuration and returns a Kubernetes clientset.
//   - Worker: A function that performs health checks on pods within a given namespace and
//     sends the results to a channel for collection and further processing. It respects
//     the context passed to it for cancellation or timeouts, and it utilizes a structured
//     logger for enhanced logging with contextual information such as the namespace and task.
//   - RunWorkers: Starts a specified number of worker goroutines that call the Worker function
//     with a given namespace and a structured logger. It returns a channel for results and a
//     function to initiate a graceful shutdown of the workers.
//
// # Usage
//
// To use this package, first create a Kubernetes client by calling NewKubernetesClient.
// Then, use the client to run worker goroutines with RunWorkers, specifying the number
// of workers, the namespace you want to target, and a context for cancellation.
//
// # Example
//
//	clientset, err := workerk8s.NewKubernetesClient()
//	if err != nil {
//	    // Handle error
//	}
//	namespace := "default" // Replace with your namespace
//	ctx := context.Background() // Use context to control worker lifetimes
//	results, shutdown := workerk8s.RunWorkers(ctx, clientset, namespace, 5)
//
//	// Do other work, then initiate graceful shutdown when needed.
//	shutdown()
//
//	// Process results until the results channel is closed.
//	for result := range results {
//	    fmt.Println(result)
//	}
//
// # Enhancements
//
//   - The Worker function now includes structured logging, which improves traceability and
//     debugging by providing contextual information in the log entries.
//   - The logging within the Worker function is now customizable, allowing different workers
//     to log with their specific contextual information such as worker index and namespace.
//
// # TODO
//
//   - Implement error handling and retry logic within the Worker function to handle transient errors.
//   - Enhance the Worker function to perform a more specific task or to be more configurable.
//   - Expand the package to support other Kubernetes resources and operations.
//   - Introduce metrics collection for monitoring the health and performance of the workers.
//
// Copyright (c) 2023 by H0llyW00dzZ
package workerk8s
