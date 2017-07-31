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

package v1alpha1

import (
	v1alpha1 "github.com/munnerz/k8s-api-pager-demo/pkg/apis/pager/v1alpha1"
	scheme "github.com/munnerz/k8s-api-pager-demo/pkg/client/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AlertsGetter has a method to return a AlertInterface.
// A group's client should implement this interface.
type AlertsGetter interface {
	Alerts(namespace string) AlertInterface
}

// AlertInterface has methods to work with Alert resources.
type AlertInterface interface {
	Create(*v1alpha1.Alert) (*v1alpha1.Alert, error)
	Update(*v1alpha1.Alert) (*v1alpha1.Alert, error)
	UpdateStatus(*v1alpha1.Alert) (*v1alpha1.Alert, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Alert, error)
	List(opts v1.ListOptions) (*v1alpha1.AlertList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Alert, err error)
	AlertExpansion
}

// alerts implements AlertInterface
type alerts struct {
	client rest.Interface
	ns     string
}

// newAlerts returns a Alerts
func newAlerts(c *PagerV1alpha1Client, namespace string) *alerts {
	return &alerts{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a alert and creates it.  Returns the server's representation of the alert, and an error, if there is any.
func (c *alerts) Create(alert *v1alpha1.Alert) (result *v1alpha1.Alert, err error) {
	result = &v1alpha1.Alert{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("alerts").
		Body(alert).
		Do().
		Into(result)
	return
}

// Update takes the representation of a alert and updates it. Returns the server's representation of the alert, and an error, if there is any.
func (c *alerts) Update(alert *v1alpha1.Alert) (result *v1alpha1.Alert, err error) {
	result = &v1alpha1.Alert{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("alerts").
		Name(alert.Name).
		Body(alert).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclientstatus=false comment above the type to avoid generating UpdateStatus().

func (c *alerts) UpdateStatus(alert *v1alpha1.Alert) (result *v1alpha1.Alert, err error) {
	result = &v1alpha1.Alert{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("alerts").
		Name(alert.Name).
		SubResource("status").
		Body(alert).
		Do().
		Into(result)
	return
}

// Delete takes name of the alert and deletes it. Returns an error if one occurs.
func (c *alerts) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("alerts").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *alerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("alerts").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the alert, and returns the corresponding alert object, and an error if there is any.
func (c *alerts) Get(name string, options v1.GetOptions) (result *v1alpha1.Alert, err error) {
	result = &v1alpha1.Alert{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("alerts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Alerts that match those selectors.
func (c *alerts) List(opts v1.ListOptions) (result *v1alpha1.AlertList, err error) {
	result = &v1alpha1.AlertList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("alerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested alerts.
func (c *alerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("alerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Patch applies the patch and returns the patched alert.
func (c *alerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Alert, err error) {
	result = &v1alpha1.Alert{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("alerts").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
