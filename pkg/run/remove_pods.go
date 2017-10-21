package run

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"

	client "github.com/srossross/k8s-test-controller/pkg/client"
	factory "github.com/srossross/k8s-test-controller/pkg/informers/externalversions"
)

// TestRunnerRemovePodsForDeletedTest will remove pods with test as owner
func TestRunnerRemovePodsForDeletedTest(
	sharedFactory factory.SharedInformerFactory,
	cl client.Interface,
	key string,
) error {
	log.Printf("  | Delete pods for removed test runner")
	Namespace, Name := splitOnce(key, "/")

	pods, err := GetPodLister(sharedFactory).Pods(Namespace).List(labels.Everything())
	if err != nil {
		return fmt.Errorf("Error getting list of pods: %s", err.Error())
	}

	pods = testRunFilter(pods, Name)
	log.Printf("  | Found %v pods to delete in namespace %s", len(pods), Namespace)

	for _, pod := range pods {
		log.Printf("  | Removing pod '%s/%s'", pod.Namespace, pod.Name)
		err := cl.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("Error removing pod: %s", err.Error())
		}
	}
	return nil
}
