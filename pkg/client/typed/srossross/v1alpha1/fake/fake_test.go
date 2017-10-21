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

package fake

import (
	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/apis/pager/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTests implements TestInterface
type FakeTests struct {
	Fake *FakeSrossrossV1alpha1
	ns   string
}

var testsResource = schema.GroupVersionResource{Group: "srossross.github.io", Version: "v1alpha1", Resource: "tests"}

var testsKind = schema.GroupVersionKind{Group: "srossross.github.io", Version: "v1alpha1", Kind: "Test"}

func (c *FakeTests) Create(test *v1alpha1.Test) (result *v1alpha1.Test, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(testsResource, c.ns, test), &v1alpha1.Test{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Test), err
}

func (c *FakeTests) Update(test *v1alpha1.Test) (result *v1alpha1.Test, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(testsResource, c.ns, test), &v1alpha1.Test{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Test), err
}

func (c *FakeTests) UpdateStatus(test *v1alpha1.Test) (*v1alpha1.Test, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(testsResource, "status", c.ns, test), &v1alpha1.Test{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Test), err
}

func (c *FakeTests) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(testsResource, c.ns, name), &v1alpha1.Test{})

	return err
}

func (c *FakeTests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(testsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.TestList{})
	return err
}

func (c *FakeTests) Get(name string, options v1.GetOptions) (result *v1alpha1.Test, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(testsResource, c.ns, name), &v1alpha1.Test{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Test), err
}

func (c *FakeTests) List(opts v1.ListOptions) (result *v1alpha1.TestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(testsResource, testsKind, c.ns, opts), &v1alpha1.TestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TestList{}
	for _, item := range obj.(*v1alpha1.TestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested tests.
func (c *FakeTests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(testsResource, c.ns, opts))

}

// Patch applies the patch and returns the patched test.
func (c *FakeTests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Test, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(testsResource, c.ns, name, data, subresources...), &v1alpha1.Test{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Test), err
}
