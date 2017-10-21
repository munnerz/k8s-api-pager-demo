package main

import (
	"flag"
	"log"
	"time"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	rest "k8s.io/client-go/rest"
	cache "k8s.io/client-go/tools/cache"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	workqueue "k8s.io/client-go/util/workqueue"

	client "github.com/srossross/k8s-test-controller/pkg/client"
	controller "github.com/srossross/k8s-test-controller/pkg/controller"
	factory "github.com/srossross/k8s-test-controller/pkg/informers/externalversions"
	run "github.com/srossross/k8s-test-controller/pkg/run"
)

var (

	// Version of this program (injected from linkflags)
	Version string

	// BuildTime of this program (injected from linkflags)
	BuildTime string

	// apiserverURL is the URL of the API server to connect to
	kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")
	// pushbulletToken is the pushbullet API token to use
	// pushbulletToken = flag.String("pushbullet-token", "", "the api token to use to send pushbullet messages")

	// queue is a queue of resources to be processed. It performs exponential
	// backoff rate limiting, with a minimum retry period of 5 seconds and a
	// maximum of 1 minute.
	rateLimiter = workqueue.NewItemExponentialFailureRateLimiter(time.Second*5, time.Minute)
	queue       = workqueue.NewRateLimitingQueue(rateLimiter)

	config *rest.Config
	// stopCh can be used to stop all the informer, as well as control loops
	// within the application.
	stopCh = make(chan struct{})

	// sharedFactory is a shared informer factory that is used a a cache for
	// items in the API server. It saves each informer listing and watching the
	// same resources independently of each other, thus providing more up to
	// date results with less 'effort'
	sharedFactory factory.SharedInformerFactory

	ctrl controller.Interface

	// cl is a Kubernetes API client for our custom resource definition type
	cl           *client.Clientset
	coreV1Client *typedv1.CoreV1Client

	// pb is the pushbullet client to use to send alerts
	// pb *pushbullet.Pushbullet
)

func main() {
	flag.Parse()

	// TODO: add proper linker flags
	log.Printf("Test controller version: %s", Version)
	log.Printf("               Built on: %s", BuildTime)

	var err error

	config, err = GetClientConfig(*kubeconfig)

	if err != nil {
		log.Fatalf("error creating config: %s", err.Error())
	}

	apiextensionsclientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		log.Fatalf("error creating api client: %s", err.Error())
	}

	err = run.InstallAllCRDs(apiextensionsclientset)

	if err != nil {
		log.Fatalf("error creating crds: %s", err.Error())
	}

	// create an instance of our own API client
	cl, err = client.NewForConfig(config)

	if err != nil {
		log.Fatalf("error creating api client: %s", err.Error())
	}

	coreV1Client, err = typedv1.NewForConfig(config)

	if err != nil {
		log.Fatalf("error creating api client: %s", err.Error())
	}

	ctrl = controller.NewTestController(&sharedFactory, cl, coreV1Client)

	log.Printf("Created Kubernetes client.")

	// we use a shared informer from the informer factory, to save calls to the
	// API as we grow our application and so state is consistent between our
	// control loops. We set a resync period of 30 seconds, in case any
	// create/replace/update/delete operations are missed when watching
	sharedFactory = factory.NewSharedInformerFactory(cl, time.Second*30)

	testRunInformer := run.NewTestRunInformer(sharedFactory, queue)

	testInformer := run.NewTestInformer(sharedFactory, queue)

	podInformer := run.SetupPodInformer(ctrl.PodInformer(), queue)

	// start the informer. This will cause it to begin receiving updates from
	// the configured API server and firing event handlers in response.
	sharedFactory.Start(stopCh)
	log.Printf("Started informer factory.")

	// wait for the informe rcache to finish performing it's initial sync of
	// resources
	if !cache.WaitForCacheSync(stopCh, testRunInformer.HasSynced) {
		log.Fatalf("error waiting for testRunInformer cache to sync: %s", err.Error())
	}

	if !cache.WaitForCacheSync(stopCh, testInformer.HasSynced) {
		log.Fatalf("error waiting for testInformer cache to sync: %s", err.Error())
	}

	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
		log.Fatalf("error waiting for podInformer cache to sync: %s", err.Error())
	}

	log.Printf("Finished populating shared informer cache.")
	// here we start just one worker reading objects off the queue. If you
	// wanted to parallelize this, you could start many instances of the worker
	// function, then ensure your application handles concurrency correctly.
	run.Work(ctrl, stopCh, queue)
}

// GetClientConfig gets config from command line kubeconfig param or InClusterConfig
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
