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

// TestRunsGetter has a method to return a TestRunInterface.
// A group's client should implement this interface.
type TestRunsGetter interface {
	TestRuns(namespace string) TestRunInterface
}

// TestRunInterface has methods to work with TestRun resources.
type TestRunInterface interface {
	Create(*v1alpha1.TestRun) (*v1alpha1.TestRun, error)
	Update(*v1alpha1.TestRun) (*v1alpha1.TestRun, error)
	UpdateStatus(*v1alpha1.TestRun) (*v1alpha1.TestRun, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.TestRun, error)
	List(opts v1.ListOptions) (*v1alpha1.TestRunList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TestRun, err error)
	TestRunExpansion
}

// testRuns implements TestRunInterface
type testRuns struct {
	client rest.Interface
	ns     string
}

// newTestRuns returns a TestRuns
func newTestRuns(c *SrossrossV1alpha1Client, namespace string) *testRuns {
	return &testRuns{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a testRun and creates it.  Returns the server's representation of the testRun, and an error, if there is any.
func (c *testRuns) Create(testRun *v1alpha1.TestRun) (result *v1alpha1.TestRun, err error) {
	result = &v1alpha1.TestRun{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("testruns").
		Body(testRun).
		Do().
		Into(result)
	return
}

// Update takes the representation of a testRun and updates it. Returns the server's representation of the testRun, and an error, if there is any.
func (c *testRuns) Update(testRun *v1alpha1.TestRun) (result *v1alpha1.TestRun, err error) {
	result = &v1alpha1.TestRun{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("testruns").
		Name(testRun.Name).
		Body(testRun).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclientstatus=false comment above the type to avoid generating UpdateStatus().

func (c *testRuns) UpdateStatus(testRun *v1alpha1.TestRun) (result *v1alpha1.TestRun, err error) {
	result = &v1alpha1.TestRun{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("testruns").
		Name(testRun.Name).
		SubResource("status").
		Body(testRun).
		Do().
		Into(result)
	return
}

// Delete takes name of the testRun and deletes it. Returns an error if one occurs.
func (c *testRuns) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("testruns").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *testRuns) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("testruns").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the testRun, and returns the corresponding testRun object, and an error if there is any.
func (c *testRuns) Get(name string, options v1.GetOptions) (result *v1alpha1.TestRun, err error) {
	result = &v1alpha1.TestRun{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("testruns").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TestRuns that match those selectors.
func (c *testRuns) List(opts v1.ListOptions) (result *v1alpha1.TestRunList, err error) {
	result = &v1alpha1.TestRunList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("testruns").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested testRuns.
func (c *testRuns) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("testruns").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Patch applies the patch and returns the patched testRun.
func (c *testRuns) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TestRun, err error) {
	result = &v1alpha1.TestRun{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("testruns").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
