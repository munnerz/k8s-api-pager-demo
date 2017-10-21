package run

import (
	// "flag"
	// "fmt"
	// "log"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgRuntime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	listerv1 "k8s.io/client-go/listers/core/v1"

	"k8s.io/client-go/tools/cache"

	"github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	"github.com/srossross/k8s-test-controller/pkg/client"
	factory "github.com/srossross/k8s-test-controller/pkg/informers/externalversions"
)

func newPodInformer(client client.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	sharedIndexInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (pkgRuntime.Object, error) {
				return client.CoreV1().Pods(metav1.NamespaceAll).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Pods(metav1.NamespaceAll).Watch(options)
			},
		},
		&corev1.Pod{},
		resyncPeriod,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	return sharedIndexInformer
}

// GetPodInformer returns an informer for a pod and registers it with the SharedInformerFactory
func GetPodInformer(sharedFactory factory.SharedInformerFactory) cache.SharedIndexInformer {
	return sharedFactory.InformerFor(&corev1.Pod{}, newPodInformer)
}

// GetPodLister returns a lister and registers it with the SharedInformerFactory
func GetPodLister(sharedFactory factory.SharedInformerFactory) listerv1.PodLister {
	return listerv1.NewPodLister(GetPodInformer(sharedFactory).GetIndexer())
}

// PodStateChange records an event for a test state change
func PodStateChange(sharedFactory factory.SharedInformerFactory, cl client.Interface, testRun *v1alpha1.TestRun, pod *corev1.Pod) error {
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
		cl, testRun, testName,
		Reason,
		fmt.Sprintf("Test pod '%s' exited with status '%s'", pod.Name, pod.Status.Phase),
	)

}
