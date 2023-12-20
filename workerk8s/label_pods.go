package workerk8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// LabelPods sets a label on all pods in a given namespace.
func LabelPods(clientset *kubernetes.Clientset, namespace, labelKey, labelValue string) error {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf(ErrorListingPods, err)
	}

	for _, pod := range pods.Items {
		if pod.Labels[labelKey] == labelValue {
			continue
		}
		podCopy := pod.DeepCopy()
		if podCopy.Labels == nil {
			podCopy.Labels = map[string]string{}
		}
		podCopy.Labels[labelKey] = labelValue
		if _, err := clientset.CoreV1().Pods(namespace).Update(context.Background(), podCopy, v1.UpdateOptions{}); err != nil {
			return fmt.Errorf(ErrorUpdatingPodLabels, err)
		}
	}
	return nil
}
