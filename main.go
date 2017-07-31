package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/munnerz/k8s-api-pager-demo/pkg/apis/pager/v1alpha1"
	"github.com/munnerz/k8s-api-pager-demo/pkg/client"
	factory "github.com/munnerz/k8s-api-pager-demo/pkg/informers/externalversions"
)

var (
	// apiserverURL is the URL of the API server to connect to
	apiserverURL = flag.String("apiserver", "http://127.0.0.1:8001", "URL used to access the Kubernetes API server")
	// pushbulletToken is the pushbullet API token to use
	pushbulletToken = flag.String("pushbullet-token", "", "the api token to use to send pushbullet messages")

	// queue is a queue of resources to be processed. It performs exponential
	// backoff rate limiting, with a minimum retry period of 5 seconds and a
	// maximum of 1 minute.
	queue = workqueue.NewRateLimitingQueue(workqueue.NewItemExponentialFailureRateLimiter(time.Second*5, time.Minute))

	// stopCh can be used to stop all the informer, as well as control loops
	// within the application.
	stopCh = make(chan struct{})

	// sharedFactory is a shared informer factory that is used a a cache for
	// items in the API server. It saves each informer listing and watching the
	// same resources independently of each other, thus providing more up to
	// date results with less 'effort'
	sharedFactory factory.SharedInformerFactory

	// cl is a Kubernetes API client for our custom resource definition type
	cl client.Interface

	// pb is the pushbullet client to use to send alerts
	pb *pushbullet.Pushbullet
)

func main() {
	flag.Parse()

	// create a client that can be used to send pushbullet notes
	pb = pushbullet.New(*pushbulletToken)

	log.Printf("Created pushbullet client.")

	var err error
	// create an instance of our own API client
	cl, err = client.NewForConfig(&rest.Config{
		Host: *apiserverURL,
	})

	if err != nil {
		log.Fatalf("error creating api client: %s", err.Error())
	}

	log.Printf("Created Kubernetes client.")

	// we use a shared informer from the informer factory, to save calls to the
	// API as we grow our application and so state is consistent between our
	// control loops. We set a resync period of 30 seconds, in case any
	// create/replace/update/delete operations are missed when watching
	sharedFactory = factory.NewSharedInformerFactory(cl, time.Second*30)

	informer := sharedFactory.Pager().V1alpha1().Alerts().Informer()
	// we add a new event handler, watching for changes to API resources.
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: enqueue,
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					enqueue(cur)
				}
			},
			DeleteFunc: enqueue,
		},
	)

	// start the informer. This will cause it to begin receiving updates from
	// the configured API server and firing event handlers in response.
	sharedFactory.Start(stopCh)
	log.Printf("Started informer factory.")

	// wait for the informe rcache to finish performing it's initial sync of
	// resources
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		log.Fatalf("error waiting for informer cache to sync: %s", err.Error())
	}

	log.Printf("Finished populating shared informer cache.")
	// here we start just one worker reading objects off the queue. If you
	// wanted to parallelize this, you could start many instances of the worker
	// function, then ensure your application handles concurrency correctly.
	work()
}

// sync will attempt to 'Sync' an alert resource. It checks to see if the alert
// has already been sent, and if not will send it and update the resource
// accordingly. This method is called whenever this controller starts, and
// whenever the resource changes, and also periodically every resyncPeriod.
func sync(al *v1alpha1.Alert) error {
	// If this message has already been sent, we exit with no error
	if al.Status.Sent {
		log.Printf("Skipping already Sent alert '%s/%s'", al.Namespace, al.Name)
		return nil
	}

	// create our note instance
	note := requests.NewNote()
	note.Title = fmt.Sprintf("Kubernetes alert for %s/%s", al.Namespace, al.Name)
	note.Body = al.Spec.Message

	// send the note. If an error occurs here, we return an error which will
	// cause the calling function to re-queue the item to be tried again later.
	if _, err := pb.PostPushesNote(note); err != nil {
		return fmt.Errorf("error sending pushbullet message: %s", err.Error())
	}
	log.Printf("Sent pushbullet note!")

	// as we've sent the note, we will update the resource accordingly.
	// if this request fails, this item will be requeued and a second alert
	// will be sent. It's therefore worth noting that this control loop will
	// send you *at least one* alert, and not *at most one*.
	al.Status.Sent = true
	if _, err := cl.PagerV1alpha1().Alerts(al.Namespace).Update(al); err != nil {
		return fmt.Errorf("error saving update to pager Alert resource: %s", err.Error())
	}
	log.Printf("Finished saving update to pager Alert resource '%s/%s'", al.Namespace, al.Name)

	// we didn't encounter any errors, so we return nil to allow the callee
	// to 'forget' this item from the queue altogether.
	return nil
}

func work() {
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

		// we define a function here to process a queue item, so that we can
		// use 'defer' to make sure the message is marked as Done on the queue
		func(key string) {
			defer queue.Done(key)

			// attempt to split the 'key' into namespace and object name
			namespace, name, err := cache.SplitMetaNamespaceKey(strKey)

			if err != nil {
				runtime.HandleError(fmt.Errorf("error splitting meta namespace key into parts: %s", err.Error()))
				return
			}

			log.Printf("Read item '%s/%s' off workqueue. Processing...", namespace, name)

			// retrieve the latest version in the cache of this alert
			obj, err := sharedFactory.Pager().V1alpha1().Alerts().Lister().Alerts(namespace).Get(name)

			if err != nil {
				runtime.HandleError(fmt.Errorf("error getting object '%s/%s' from api: %s", namespace, name, err.Error()))
				return
			}

			log.Printf("Got most up to date version of '%s/%s'. Syncing...", namespace, name)

			// attempt to sync the current state of the world with the desired!
			// If sync returns an error, we skip calling `queue.Forget`,
			// thus causing the resource to be requeued at a later time.
			if err := sync(obj); err != nil {
				runtime.HandleError(fmt.Errorf("error processing item '%s/%s': %s", namespace, name, err.Error()))
				return
			}

			log.Printf("Finished processing '%s/%s' successfully! Removing from queue.", namespace, name)

			// as we managed to process this successfully, we can forget it
			// from the work queue altogether.
			queue.Forget(key)
		}(strKey)
	}
}

// enqueue will add an object 'obj' into the workqueue. The object being added
// must be of type metav1.Object, metav1.ObjectAccessor or cache.ExplicitKey.
func enqueue(obj interface{}) {
	// DeletionHandlingMetaNamespaceKeyFunc will convert an object into a
	// 'namespace/name' string. We do this because our item may be processed
	// much later than now, and so we want to ensure it gets a fresh copy of
	// the resource when it starts. Also, this allows us to keep adding the
	// same item into the work queue without duplicates building up.
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("error obtaining key for object being enqueue: %s", err.Error()))
		return
	}
	// add the item to the queue
	queue.Add(key)
}
