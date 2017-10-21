package run

import (
	"fmt"
	"log"
	"strings"

	corev1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/util/runtime"
	workqueue "k8s.io/client-go/util/workqueue"

	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	controller "github.com/srossross/k8s-test-controller/pkg/controller"
)

func splitOnce(key, sep string) (string, string) {
	tmp := strings.SplitN(key, sep, 2)
	if len(tmp) == 1 {
		return tmp[0], ""
	}
	return tmp[0], tmp[1]
}

// Work pops jobs off the queue
func Work(ctrl controller.Interface, stopCh chan struct{}, queue workqueue.RateLimitingInterface) {
	for {
		// we read a message off the queue
		key, shutdown := queue.Get()

		// if the queue has been shut down, we should exit the work queue here
		if shutdown {
			stopCh <- struct{}{}
			return
		}

		// convert the queue item into a string. If it's not a string, we'll
		// simply discard it as invalid data and log a message.
		var strKey string
		var ok bool
		if strKey, ok = key.(string); !ok {
			runtime.HandleError(fmt.Errorf("key in queue should be of type string but got %T. discarding", key))
			return
		}

		log.Printf("Popped '%s' off the queue", key)
		// we define a function here to process a queue item, so that we can
		// use 'defer' to make sure the message is marked as Done on the queue
		func(key string) {
			defer queue.Done(key)
			runType, key := splitOnce(key, ":")

			var err error
			var testRun *v1alpha1.TestRun
			var pod *corev1.Pod

			switch runType {
			case ReconsilePodStatus:
				{
					testRun, pod, err = ctrl.GetPodAndTestRunFromKey(key)
					if err == nil {
						err = PodStateChange(ctrl, testRun, pod)
					}
				}
			case ReconsileTestRun:
				{
					testRun, err = ctrl.GetTestRunFromKey(key)
					if err == nil {
						err = TestRunner(ctrl, testRun)
					} else if errors.IsNotFound(err) {
						// FIXME: should this be handled by k8s garbage collection?
						err = ctrl.TestRunnerRemovePodsForDeletedTest(key)
					}
				}
			default:
				err = fmt.Errorf("key in queue should be of type string but got %T. discarding", key)
			}

			if err != nil {
				runtime.HandleError(err)
				return
			}
			queue.Forget(key)
		}(strKey)
	}
}
