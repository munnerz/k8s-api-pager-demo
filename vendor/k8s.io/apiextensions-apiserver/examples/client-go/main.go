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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	crv1 "k8s.io/apiextensions-apiserver/examples/client-go/apis/cr/v1"
	exampleclient "k8s.io/apiextensions-apiserver/examples/client-go/client"
	examplecontroller "k8s.io/apiextensions-apiserver/examples/client-go/controller"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}


	apiextensionsclientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// initialize custom resource using a CustomResourceDefinition if it does not exist
	// crd, err := exampleclient.CreateCustomResourceDefinition(apiextensionsclientset)
	// if err != nil && !apierrors.IsAlreadyExists(err) {
	// 	panic(err)
	// }
	//
	// if crd != nil {
	// 	defer apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)
	// }
	crdName := crv1.E2ETestResourcePlural + "." + crv1.GroupName
	crd, err := apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("CRD %v", crd)

	// make a new config for our extension's API group, using the first config as a baseline
	exampleClient, exampleScheme, err := exampleclient.NewClient(config)
	if err != nil {
		panic(err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	crv1.NewE2ETestController(client, *exampleClient).Run(stopCh)

	// start a controller on instances of our custom resource
	controller := examplecontroller.E2ETestController{
		E2ETestClient: exampleClient,
		E2ETestScheme: exampleScheme,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go controller.Run(ctx)

	// Create an instance of our custom resource
	example := &crv1.E2ETest{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test1",
			// Namespace: "default",
		},
		Spec: crv1.E2ETestSpec{
			Foo: "hello",
			Bar: true,
		},
		Status: crv1.E2ETestStatus{
			State:   crv1.E2ETestStateCreated,
			Message: "Created, not processed yet",
		},
	}
	var result crv1.E2ETest
	fmt.Printf("WANT: %#v\n", example)
	err = exampleClient.Post().
		Resource(crv1.E2ETestResourcePlural).
		Namespace(apiv1.NamespaceDefault).
		Body(example).
		Do().Into(&result)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}

	// Poll until E2ETest object is handled by controller and gets status updated to "Processed"
	err = exampleclient.WaitForE2ETestInstanceProcessed(exampleClient, "test1")
	if err != nil {
		panic(err)
	}
	fmt.Print("PROCESSED\n")

	// Fetch a list of our CRs
	exampleList := crv1.E2ETestList{}
	err = exampleClient.Get().Resource(crv1.E2ETestResourcePlural).Do().Into(&exampleList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleList)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
