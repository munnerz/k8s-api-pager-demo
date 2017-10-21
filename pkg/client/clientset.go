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

package client

import (
	glog "github.com/golang/glog"
	srossrossv1alpha1 "github.com/srossross/k8s-test-controller/pkg/client/typed/srossross/v1alpha1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	SrossrossV1alpha1() srossrossv1alpha1.SrossrossV1alpha1Interface
	CoreV1() corev1.CoreV1Interface
	// Deprecated: please explicitly pick a version if possible.
	Srossross() srossrossv1alpha1.SrossrossV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	*srossrossv1alpha1.SrossrossV1alpha1Client
	*corev1.CoreV1Client
}

// SrossrossV1alpha1 retrieves the SrossrossV1alpha1Client
func (c *Clientset) SrossrossV1alpha1() srossrossv1alpha1.SrossrossV1alpha1Interface {
	if c == nil {
		return nil
	}
	return c.SrossrossV1alpha1Client
}

// Deprecated: Srossross retrieves the default version of SrossrossClient.
// Please explicitly pick a version.
func (c *Clientset) Srossross() srossrossv1alpha1.SrossrossV1alpha1Interface {
	if c == nil {
		return nil
	}
	return c.SrossrossV1alpha1Client
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// CoreV1 retrieves the CoreV1Client
func (c *Clientset) CoreV1() corev1.CoreV1Interface {
	if c == nil {
		return nil
	}
	return c.CoreV1Client
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.SrossrossV1alpha1Client, err = srossrossv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}

	cs.CoreV1Client, err = corev1.NewForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the CoreV1Client: %v", err)
		return nil, err
	}

	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.SrossrossV1alpha1Client = srossrossv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.SrossrossV1alpha1Client = srossrossv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
