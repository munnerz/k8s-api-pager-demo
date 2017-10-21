package controller

import (
	"k8s.io/api/core/v1"
)

// TestRunFilter filter pods that were instantiated with testrun of name
// testRunName
func TestRunFilter(pods []*v1.Pod, testRunName string) []*v1.Pod {
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
