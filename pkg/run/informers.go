package run

import (
	"fmt"
	"reflect"

	log "github.com/golang/glog"

	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	factory "github.com/srossross/k8s-test-controller/pkg/informers/externalversions"
	"k8s.io/api/core/v1"
	cache "k8s.io/client-go/tools/cache"
	workqueue "k8s.io/client-go/util/workqueue"
)

type ReconsileType string

var (
	ReconsilePodStatus = "Pod"
	ReconsileTestRun   = "TestRun"
)

func isStatusChange(old, cur interface{}) bool {
	oldPod, ok := old.(*v1.Pod)
	if !ok {
		return false
	}
	curPod, ok := cur.(*v1.Pod)
	if !ok {
		return false
	}

	if oldPod.Status.Phase != curPod.Status.Phase {
		log.Infof(
			"Pod '%v/%v' phase changed from '%s' to '%s'",
			curPod.Namespace, curPod.Name,
			oldPod.Status.Phase, curPod.Status.Phase,
		)
		return true
	}
	return false
}

func testRunKey(cur interface{}) (string, bool) {

	testRun, ok := cur.(*v1alpha1.TestRun)

	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v:%v/%v", ReconsileTestRun, testRun.Namespace, testRun.Name), true
}

func podTestRunKey(cur interface{}) (string, bool) {
	pod, ok := cur.(*v1.Pod)
	if !ok {
		return "", false
	}
	annotaion, ok := pod.Annotations["srossross.github.io/v1alpha1"]
	if !ok {
		return "", false
	}
	return annotaion, true
}

// NewTestRunInformer creates a new test run Informer that watches and caches testruns
func NewTestRunInformer(
	sharedFactory factory.SharedInformerFactory,
	queue workqueue.RateLimitingInterface,
) cache.SharedIndexInformer {

	runInformer := sharedFactory.Srossross().V1alpha1().TestRuns().Informer()
	// we add a new event handler, watching for changes to API resources.

	enqueue := func(cur interface{}) {
		key, ok := testRunKey(cur)
		if !ok {
			log.Infof("Error getting testrun queue key")
			return
		}
		queue.Add(key)
	}

	runInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(cur interface{}) { enqueue(cur) },
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					enqueue(cur)
				}
			},
			DeleteFunc: func(cur interface{}) { enqueue(cur) },
		},
	)

	return runInformer
}

// NewTestInformer creates a new test Informer that watches and caches tests
func NewTestInformer(sharedFactory factory.SharedInformerFactory, queue workqueue.RateLimitingInterface) cache.SharedIndexInformer {
	testInformer := sharedFactory.Srossross().V1alpha1().Tests().Informer()
	// we add a new event handler, watching for changes to API resources.
	testInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(cur interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cur)
				if err != nil {
					log.Fatalf("Error in DeletionHandlingMetaNamespaceKeyFunc %v", err.Error())
				}
				log.V(4).Infof("Test %v Added (not triggering reconsile loop)", key)
			},
			UpdateFunc: func(old, cur interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cur)
				if err != nil {
					log.Fatalf("Error in DeletionHandlingMetaNamespaceKeyFunc %v", err.Error())
				}
				log.V(4).Infof("Test %v Updated (not triggering reconsile loop)", key)
			},
			DeleteFunc: func(cur interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cur)
				if err != nil {
					log.Fatalf("Error in DeletionHandlingMetaNamespaceKeyFunc %v", err.Error())
				}
				log.V(4).Infof("Test %v Deleted (not triggering reconsile loop)", key)
			},
		},
	)
	return testInformer
}

// SetupPodInformer  creates a new test Informer that watches and caches pods
func SetupPodInformer(podInformer cache.SharedIndexInformer, queue workqueue.RateLimitingInterface) cache.SharedIndexInformer {

	enqueue := func(cur interface{}) {
		key, ok := podTestRunKey(cur)
		if !ok {
			// log.Infof("Error getting testrun queue key")
			return
		}
		queue.Add(key)
	}

	enqueuePodStatEvent := func(cur interface{}) {
		_, ok := podTestRunKey(cur)
		if !ok {
			return
		}
		key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cur)
		if err != nil {
			return
		}
		queue.Add(fmt.Sprintf("%s:%s", ReconsilePodStatus, key))
	}

	podInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(cur interface{}) { enqueue(cur) },
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					// FIXME: we should detect a change in state so that
					// we can add an test fail/success event
					if isStatusChange(old, cur) {
						enqueuePodStatEvent(cur)
					}
					enqueue(cur)
				}
			},
			DeleteFunc: func(cur interface{}) { enqueue(cur) },
		},
	)
	return podInformer
}
