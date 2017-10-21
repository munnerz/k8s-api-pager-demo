package controller

import (
	corev1 "k8s.io/api/core/v1"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"

	"github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	client "github.com/srossross/k8s-test-controller/pkg/client"
	srossrossv1alpha1 "github.com/srossross/k8s-test-controller/pkg/client/typed/srossross/v1alpha1"
	factory "github.com/srossross/k8s-test-controller/pkg/informers/externalversions"
	listerV1alpha1 "github.com/srossross/k8s-test-controller/pkg/listers/srossross/v1alpha1"
)

// Interface for a TestController
type Interface interface {
	TestLister() listerV1alpha1.TestLister
	TestRunLister() listerV1alpha1.TestRunLister
	SrossrossV1alpha1() srossrossv1alpha1.SrossrossV1alpha1Interface

	PodLister() listerv1.PodLister
	PodInformer() cache.SharedIndexInformer
	CoreV1() typedv1.CoreV1Interface

	GetTestRunFromKey(key string) (*v1alpha1.TestRun, error)
	GetPodAndTestRunFromKey(key string) (*v1alpha1.TestRun, *corev1.Pod, error)
	TestRunnerRemovePodsForDeletedTest(key string) error
}

// TestController creates a single interface to run the reconsile loop
type TestController struct {
	sharedFactory *factory.SharedInformerFactory
	testClient    *client.Clientset
	coreV1Client  *typedv1.CoreV1Client
}

// CoreV1 get CoreV1 client
func (ctrl *TestController) CoreV1() typedv1.CoreV1Interface {
	if ctrl == nil {
		return nil
	}
	return ctrl.coreV1Client
}

// SrossrossV1alpha1 get SrossrossV1alpha1 client
func (ctrl *TestController) SrossrossV1alpha1() srossrossv1alpha1.SrossrossV1alpha1Interface {
	return (*ctrl.testClient).SrossrossV1alpha1()
}

// TestLister get a testlister
func (ctrl *TestController) TestLister() listerV1alpha1.TestLister {
	return (*ctrl.sharedFactory).Srossross().V1alpha1().Tests().Lister()
}

// TestRunLister get a testrun lister
func (ctrl *TestController) TestRunLister() listerV1alpha1.TestRunLister {
	return (*ctrl.sharedFactory).Srossross().V1alpha1().TestRuns().Lister()
}

// PodInformer returns an informer for a pod and registers it with the SharedInformerFactory
func (ctrl *TestController) PodInformer() cache.SharedIndexInformer {
	return (*ctrl.sharedFactory).InformerFor(&corev1.Pod{}, ctrl.newPodInformerFactory())
}

// PodLister gets a corev1 podlister
func (ctrl *TestController) PodLister() listerv1.PodLister {
	return listerv1.NewPodLister(ctrl.PodInformer().GetIndexer())
}

//NewTestController creates a new TestController
func NewTestController(sharedFactory *factory.SharedInformerFactory, testClient *client.Clientset, coreV1Client *typedv1.CoreV1Client) *TestController {
	return &TestController{
		sharedFactory: sharedFactory,
		testClient:    testClient,
		coreV1Client:  coreV1Client,
	}
}
