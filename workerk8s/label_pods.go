package workerk8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// LabelPods sets a label on all pods in a given namespace.
func LabelPods(clientset *kubernetes.Clientset, namespace string, labelKey string, labelValue string) error {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf(errorlistingspods, err)
	}

	for _, pod := range pods.Items {
		// Skip if the pod already has the label with the desired value
		if value, ok := pod.Labels[labelKey]; ok && value == labelValue {
			continue
		}

		// Deep copy the pod to avoid modifying the pod list items directly
		podToUpdate := pod.DeepCopy()

		// Ensure the labels map is not nil
		if podToUpdate.Labels == nil {
			podToUpdate.Labels = make(map[string]string)
		}

		// Set the label
		podToUpdate.Labels[labelKey] = labelValue

		// Update the pod with the new label
		_, err := clientset.CoreV1().Pods(namespace).Update(context.TODO(), podToUpdate, v1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf(errorupdatingpodlabels, err)
		}
	}

	return nil
}
