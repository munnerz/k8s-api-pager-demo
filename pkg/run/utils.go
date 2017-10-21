package run

import (
	"strings"

	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	factory "github.com/srossross/k8s-test-controller/pkg/informers/externalversions"
	"k8s.io/api/core/v1"
)

func testRunFilter(pods []*v1.Pod, testRunName string) []*v1.Pod {
	result := []*v1.Pod{}
	for _, pod := range pods {
		if pod.Labels["test-run"] == testRunName {
			result = append(result, pod)
		}
		// if(!strings.HasPrefix(a[i], "foo_") && len(a[i]) <= 7) {
		//     nofoos = append(nofoos, a[i])
	}

	return result
}

func splitOnce(key, sep string) (string, string) {
	tmp := strings.SplitN(key, sep, 2)
	if len(tmp) == 1 {
		return tmp[0], ""
	}
	return tmp[0], tmp[1]
}

// GetTestRunFromKey get a test run from a key put on the queue
func GetTestRunFromKey(sharedFactory factory.SharedInformerFactory, key string) (*v1alpha1.TestRun, error) {

	namespace, name := splitOnce(key, "/")
	lister := sharedFactory.Srossross().V1alpha1().TestRuns().Lister()
	return lister.TestRuns(namespace).Get(name)
}

// GetPodAndTestRunFromKey get a test run from a key put on the queue
func GetPodAndTestRunFromKey(sharedFactory factory.SharedInformerFactory, key string) (*v1alpha1.TestRun, *v1.Pod, error) {

	namespace, name := splitOnce(key, "/")

	pod, err := GetPodLister(sharedFactory).Pods(namespace).Get(name)
	if err != nil {
		return nil, nil, err
	}

	lister := sharedFactory.Srossross().V1alpha1().TestRuns().Lister()
	testRun, err := lister.TestRuns(namespace).Get(pod.Labels["test-run"])

	if err != nil {
		return nil, nil, err
	}

	return testRun, pod, nil

}
