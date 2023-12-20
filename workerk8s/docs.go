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
//     the context passed to it for cancellation or timeouts.
//   - RunWorkers: Starts a specified number of worker goroutines that call the Worker function
//     with a given namespace and collects the results from all workers into a slice of strings.
//
// # Usage
//
// To use this package, first create a Kubernetes client by calling NewKubernetesClient.
// Then, use the client to run worker goroutines with RunWorkers, specifying the number
// of workers and the namespace you want to target.
//
// # Example
//
//	clientset, err := workerk8s.NewKubernetesClient()
//	if err != nil {
//	    // Handle error
//	}
//	namespace := "default" // Replace with your namespace
//	results := workerk8s.RunWorkers(clientset, 5, namespace)
//	for _, result := range results {
//	    fmt.Println(result)
//	}
//
// # TODO
//
//   - Implement error handling and retry logic within the Worker function to handle transient errors.
//   - Enhance the Worker function to perform a more specific task or to be more configurable.
//   - Consider adding a function to clean up resources or to gracefully shut down the workers.
//   - Expand the package to support other Kubernetes resources and operations.
//   - Add context parameter to RunWorkers to control the lifetime of worker operations and allow for shutdown signals.
//
// Copyright (c) 2023 by H0llyW00dzZ
package workerk8s
