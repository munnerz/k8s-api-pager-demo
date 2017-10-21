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

// FakeTestRuns implements TestRunInterface
type FakeTestRuns struct {
	Fake *FakeSrossrossV1alpha1
	ns   string
}

var testrunsResource = schema.GroupVersionResource{Group: "srossross.github.io", Version: "v1alpha1", Resource: "testruns"}

var testrunsKind = schema.GroupVersionKind{Group: "srossross.github.io", Version: "v1alpha1", Kind: "TestRun"}

func (c *FakeTestRuns) Create(testRun *v1alpha1.TestRun) (result *v1alpha1.TestRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(testrunsResource, c.ns, testRun), &v1alpha1.TestRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TestRun), err
}

func (c *FakeTestRuns) Update(testRun *v1alpha1.TestRun) (result *v1alpha1.TestRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(testrunsResource, c.ns, testRun), &v1alpha1.TestRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TestRun), err
}

func (c *FakeTestRuns) UpdateStatus(testRun *v1alpha1.TestRun) (*v1alpha1.TestRun, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(testrunsResource, "status", c.ns, testRun), &v1alpha1.TestRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TestRun), err
}

func (c *FakeTestRuns) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(testrunsResource, c.ns, name), &v1alpha1.TestRun{})

	return err
}

func (c *FakeTestRuns) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(testrunsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.TestRunList{})
	return err
}

func (c *FakeTestRuns) Get(name string, options v1.GetOptions) (result *v1alpha1.TestRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(testrunsResource, c.ns, name), &v1alpha1.TestRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TestRun), err
}

func (c *FakeTestRuns) List(opts v1.ListOptions) (result *v1alpha1.TestRunList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(testrunsResource, testrunsKind, c.ns, opts), &v1alpha1.TestRunList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TestRunList{}
	for _, item := range obj.(*v1alpha1.TestRunList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested testRuns.
func (c *FakeTestRuns) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(testrunsResource, c.ns, opts))

}

// Patch applies the patch and returns the patched testRun.
func (c *FakeTestRuns) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TestRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(testrunsResource, c.ns, name, data, subresources...), &v1alpha1.TestRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TestRun), err
}
