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
	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	scheme "github.com/srossross/k8s-test-controller/pkg/client/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TestsGetter has a method to return a TestInterface.
// A group's client should implement this interface.
type TestsGetter interface {
	Tests(namespace string) TestInterface
}

// TestInterface has methods to work with Test resources.
type TestInterface interface {
	Create(*v1alpha1.Test) (*v1alpha1.Test, error)
	Update(*v1alpha1.Test) (*v1alpha1.Test, error)
	UpdateStatus(*v1alpha1.Test) (*v1alpha1.Test, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Test, error)
	List(opts v1.ListOptions) (*v1alpha1.TestList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Test, err error)
	TestExpansion
}

// tests implements TestInterface
type tests struct {
	client rest.Interface
	ns     string
}

// newTests returns a Tests
func newTests(c *SrossrossV1alpha1Client, namespace string) *tests {
	return &tests{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a test and creates it.  Returns the server's representation of the test, and an error, if there is any.
func (c *tests) Create(test *v1alpha1.Test) (result *v1alpha1.Test, err error) {
	result = &v1alpha1.Test{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("tests").
		Body(test).
		Do().
		Into(result)
	return
}

// Update takes the representation of a test and updates it. Returns the server's representation of the test, and an error, if there is any.
func (c *tests) Update(test *v1alpha1.Test) (result *v1alpha1.Test, err error) {
	result = &v1alpha1.Test{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tests").
		Name(test.Name).
		Body(test).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclientstatus=false comment above the type to avoid generating UpdateStatus().

func (c *tests) UpdateStatus(test *v1alpha1.Test) (result *v1alpha1.Test, err error) {
	result = &v1alpha1.Test{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tests").
		Name(test.Name).
		SubResource("status").
		Body(test).
		Do().
		Into(result)
	return
}

// Delete takes name of the test and deletes it. Returns an error if one occurs.
func (c *tests) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tests").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the test, and returns the corresponding test object, and an error if there is any.
func (c *tests) Get(name string, options v1.GetOptions) (result *v1alpha1.Test, err error) {
	result = &v1alpha1.Test{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Tests that match those selectors.
func (c *tests) List(opts v1.ListOptions) (result *v1alpha1.TestList, err error) {
	result = &v1alpha1.TestList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tests.
func (c *tests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("tests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Patch applies the patch and returns the patched test.
func (c *tests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Test, err error) {
	result = &v1alpha1.Test{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("tests").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
