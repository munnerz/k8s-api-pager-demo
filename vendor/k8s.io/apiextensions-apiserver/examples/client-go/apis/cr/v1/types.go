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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

const E2ETestResourcePlural = "tests"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type E2ETest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              E2ETestSpec   `json:"spec"`
	Status            E2ETestStatus `json:"status,omitempty"`
}

type E2ETestSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
	Template corev1.PodTemplateSpec `json:"template" protobuf:"bytes,3,opt,name=template"`
}



type E2ETestStatus struct {
	State   E2ETestState `json:"state,omitempty"`
	Message string       `json:"message,omitempty"`
}

type E2ETestState string

const (
	E2ETestStateCreated   E2ETestState = "Created"
	E2ETestStateProcessed E2ETestState = "Processed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type E2ETestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []E2ETest `json:"items"`
}
