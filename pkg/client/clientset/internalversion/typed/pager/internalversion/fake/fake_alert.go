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
	pager "github.com/munnerz/k8s-api-pager-demo/pkg/apis/pager"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAlerts implements AlertInterface
type FakeAlerts struct {
	Fake *FakePager
	ns   string
}

var alertsResource = schema.GroupVersionResource{Group: "pager.k8s.co", Version: "", Resource: "alerts"}

var alertsKind = schema.GroupVersionKind{Group: "pager.k8s.co", Version: "", Kind: "Alert"}

// Get takes name of the alert, and returns the corresponding alert object, and an error if there is any.
func (c *FakeAlerts) Get(name string, options v1.GetOptions) (result *pager.Alert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(alertsResource, c.ns, name), &pager.Alert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*pager.Alert), err
}

// List takes label and field selectors, and returns the list of Alerts that match those selectors.
func (c *FakeAlerts) List(opts v1.ListOptions) (result *pager.AlertList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(alertsResource, alertsKind, c.ns, opts), &pager.AlertList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &pager.AlertList{}
	for _, item := range obj.(*pager.AlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested alerts.
func (c *FakeAlerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(alertsResource, c.ns, opts))

}

// Create takes the representation of a alert and creates it.  Returns the server's representation of the alert, and an error, if there is any.
func (c *FakeAlerts) Create(alert *pager.Alert) (result *pager.Alert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(alertsResource, c.ns, alert), &pager.Alert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*pager.Alert), err
}

// Update takes the representation of a alert and updates it. Returns the server's representation of the alert, and an error, if there is any.
func (c *FakeAlerts) Update(alert *pager.Alert) (result *pager.Alert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(alertsResource, c.ns, alert), &pager.Alert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*pager.Alert), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAlerts) UpdateStatus(alert *pager.Alert) (*pager.Alert, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(alertsResource, "status", c.ns, alert), &pager.Alert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*pager.Alert), err
}

// Delete takes name of the alert and deletes it. Returns an error if one occurs.
func (c *FakeAlerts) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(alertsResource, c.ns, name), &pager.Alert{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAlerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(alertsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &pager.AlertList{})
	return err
}

// Patch applies the patch and returns the patched alert.
func (c *FakeAlerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *pager.Alert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(alertsResource, c.ns, name, data, subresources...), &pager.Alert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*pager.Alert), err
}
