package workerk8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// LabelPods sets a specific label on all pods within the specified namespace that do not already have it.
func LabelPods(clientset *kubernetes.Clientset, namespace, labelKey, labelValue string) error {
	// Retrieve a list of all pods in the given namespace.
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf(ErrorListingPods, err)
	}

	// Iterate over the list of pods and update their labels if necessary.
	for _, pod := range pods.Items {
		if err := labelSinglePod(clientset, &pod, namespace, labelKey, labelValue); err != nil {
			return err
		}
	}
	return nil
}

// labelSinglePod applies the label to a single pod if it doesn't already have it.
func labelSinglePod(clientset *kubernetes.Clientset, pod *corev1.Pod, namespace, labelKey, labelValue string) error {
	// If the pod already has the label with the correct value, skip updating.
	if pod.Labels[labelKey] == labelValue {
		return nil
	}

	// Prepare the pod's labels for update or create a new label map if none exist.
	podCopy := pod.DeepCopy()
	if podCopy.Labels == nil {
		podCopy.Labels = make(map[string]string)
	}
	podCopy.Labels[labelKey] = labelValue

	// Update the pod with the new label.
	_, err := clientset.CoreV1().Pods(namespace).Update(context.Background(), podCopy, v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf(ErrorUpdatingPodLabels, err)
	}
	return nil
}
