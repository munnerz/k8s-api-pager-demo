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
	// corev1 "k8s.io/api/core/v1"
)

const E2ETestRunnerResourcePlural = "testrunners"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type E2ETestRunner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              E2ETestRunnerSpec   `json:"spec"`
	Status            E2ETestRunnerStatus `json:"status,omitempty"`
}

type E2ETestRunnerSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}



type E2ETestRunnerStatus struct {
	State   E2ETestRunnerState `json:"state,omitempty"`
	Message string       `json:"message,omitempty"`
}

type E2ETestRunnerState string

const (
	E2ETestRunnerStateCreated   E2ETestRunnerState = "Created"
	E2ETestRunnerStateProcessed E2ETestRunnerState = "Processed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type E2ETestRunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []E2ETestRunner `json:"items"`
}
