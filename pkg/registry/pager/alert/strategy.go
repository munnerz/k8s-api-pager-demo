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

package alert

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	"github.com/munnerz/k8s-api-pager-demo/pkg/apis/pager"
)

func NewStrategy(typer runtime.ObjectTyper) alertStrategy {
	return alertStrategy{typer, names.SimpleNameGenerator}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	apiserver, ok := obj.(*pager.Alert)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not an Alert")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), AlertToSelectableFields(apiserver), apiserver.Initializers != nil, nil
}

// MatchAlert is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchAlert(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// AlertToSelectableFields returns a field set that represents the object.
func AlertToSelectableFields(obj *pager.Alert) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type alertStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (alertStrategy) NamespaceScoped() bool {
	return true
}

func (alertStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
}

func (alertStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
}

func (alertStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (alertStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (alertStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (alertStrategy) Canonicalize(obj runtime.Object) {
}

func (alertStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// implements interface RESTUpdateStrategy. This implementation validates updates to
// instance.Status updates only and disallows any modifications to the instance.Spec.
type alertStatusStrategy struct {
	alertStrategy
}

func (alertStatusStrategy) PrepareForUpdate(ctx genericapirequest.Context, new, old runtime.Object) {
	newAlert, ok := new.(*pager.Alert)
	if !ok {
		glog.Fatal("received a non-alert object to update to")
	}
	oldAlert, ok := old.(*pager.Alert)
	if !ok {
		glog.Fatal("received a non-alert object to update from")
	}
	// Status changes are not allowed to update spec
	newAlert.Spec = oldAlert.Spec
}

func (alertStatusStrategy) ValidateUpdate(ctx genericapirequest.Context, new, old runtime.Object) field.ErrorList {
	return nil
}
