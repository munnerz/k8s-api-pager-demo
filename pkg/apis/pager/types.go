package pager

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Alert struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   AlertSpec
	Status AlertStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AlertList struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Items []Alert
}

type AlertSpec struct {
	Message string
}

type AlertStatus struct {
	Sent bool
}
