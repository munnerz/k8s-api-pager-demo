package controller

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
)

// TestRunnerRemovePodsForDeletedTest will remove pods with test as owner
func (ctrl *TestController) TestRunnerRemovePodsForDeletedTest(
	key string,
) error {
	log.Printf("  | Delete pods for removed test runner")
	Namespace, Name := splitOnce(key, "/")

	pods, err := ctrl.PodLister().Pods(Namespace).List(labels.Everything())
	if err != nil {
		return fmt.Errorf("Error getting list of pods: %s", err.Error())
	}

	pods = TestRunFilter(pods, Name)
	log.Printf("  | Found %v pods to delete in namespace %s", len(pods), Namespace)

	for _, pod := range pods {
		log.Printf("  | Removing pod '%s/%s'", pod.Namespace, pod.Name)
		err := ctrl.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("Error removing pod: %s", err.Error())
		}
	}
	return nil
}
