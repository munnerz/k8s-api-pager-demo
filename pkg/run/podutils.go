package run

import (
	// "flag"
	// "fmt"
	// "log"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	controller "github.com/srossross/k8s-test-controller/pkg/controller"
)

// PodStateChange records an event for a test state change
func PodStateChange(ctrl controller.Interface, testRun *v1alpha1.TestRun, pod *corev1.Pod) error {
	testName, ok := pod.Labels["test-name"]
	if !ok {
		return fmt.Errorf("Could not get test-name label from pod %s", pod.Name)
	}
	var Reason string
	switch pod.Status.Phase {
	case "Succeeded":
		Reason = "TestSuccess"
	case "Failed":
		Reason = "TestFail"
	case "Unknown":
		Reason = "TestError"
	case "Pending":
		return nil
	case "Running":
		return nil
	}

	return CreateTestRunEvent(
		ctrl, testRun, testName,
		Reason,
		fmt.Sprintf("Test pod '%s' exited with status '%s'", pod.Name, pod.Status.Phase),
	)

}
