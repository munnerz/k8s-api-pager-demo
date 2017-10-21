package controller

import (
	"strings"

	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	"k8s.io/api/core/v1"
)

func splitOnce(key, sep string) (string, string) {
	tmp := strings.SplitN(key, sep, 2)
	if len(tmp) == 1 {
		return tmp[0], ""
	}
	return tmp[0], tmp[1]
}

// GetTestRunFromKey get a test run from a key put on the queue
func (ctrl *TestController) GetTestRunFromKey(key string) (*v1alpha1.TestRun, error) {
	namespace, name := splitOnce(key, "/")
	return ctrl.TestRunLister().TestRuns(namespace).Get(name)
}

// GetPodAndTestRunFromKey get a test run from a key put on the queue
func (ctrl *TestController) GetPodAndTestRunFromKey(key string) (*v1alpha1.TestRun, *v1.Pod, error) {

	namespace, name := splitOnce(key, "/")

	pod, err := ctrl.PodLister().Pods(namespace).Get(name)
	if err != nil {
		return nil, nil, err
	}

	testRun, err := ctrl.TestRunLister().TestRuns(namespace).Get(pod.Labels["test-run"])

	if err != nil {
		return nil, nil, err
	}

	return testRun, pod, nil

}
