/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	crv1 "k8s.io/apiextensions-apiserver/examples/client-go/apis/cr/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type E2ETestController struct {
	E2ETestClient *rest.RESTClient
	E2ETestScheme *runtime.Scheme
}

// Run starts an E2ETest resource controller
func (c *E2ETestController) Run(ctx context.Context) error {
	fmt.Print("Watch E2ETest objects\n")

	// Watch E2ETest objects
	_, err := c.watchE2ETests(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for E2ETest resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *E2ETestController) watchE2ETests(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.E2ETestClient,
		crv1.E2ETestResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&crv1.E2ETest{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *E2ETestController) onAdd(obj interface{}) {
	example := obj.(*crv1.E2ETest)
	fmt.Printf("[CONTROLLER] OnAdd %s\n", example.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	exampleCopy := example.DeepCopy()
	exampleCopy.Status = crv1.E2ETestStatus{
		State:   crv1.E2ETestStateProcessed,
		Message: "Successfully processed by controller",
	}

	err := c.E2ETestClient.Put().
		Name(example.ObjectMeta.Name).
		// Namespace(example.ObjectMeta.Namespace).
		Resource(crv1.E2ETestResourcePlural).
		Body(exampleCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", exampleCopy)
	}
}

func (c *E2ETestController) onUpdate(oldObj, newObj interface{}) {
	oldE2ETest := oldObj.(*crv1.E2ETest)
	newE2ETest := newObj.(*crv1.E2ETest)
	fmt.Printf("[CONTROLLER] OnUpdate oldObj: %s\n", oldE2ETest.ObjectMeta.SelfLink)
	fmt.Printf("[CONTROLLER] OnUpdate newObj: %s\n", newE2ETest.ObjectMeta.SelfLink)
}

func (c *E2ETestController) onDelete(obj interface{}) {
	example := obj.(*crv1.E2ETest)
	fmt.Printf("[CONTROLLER] OnDelete %s\n", example.ObjectMeta.SelfLink)
}
