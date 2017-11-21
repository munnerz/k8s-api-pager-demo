package v1alpha1

import (
	"k8s.io/apimachinery/pkg/conversion"

	"github.com/munnerz/k8s-api-pager-demo/pkg/apis/pager"
)

const (
	AnnotationCustomTitle = "pager.k8s.co/title"
)

func Convert_pager_Alert_To_v1alpha1_Alert(in *pager.Alert, out *Alert, s conversion.Scope) error {
	err := autoConvert_pager_Alert_To_v1alpha1_Alert(in, out, s)
	if err != nil {
		return err
	}
	if out.ObjectMeta.Annotations == nil {
		out.ObjectMeta.Annotations = map[string]string{}
	}
	out.ObjectMeta.Annotations[AnnotationCustomTitle] = in.Spec.Title
	return nil
}

func Convert_v1alpha1_Alert_To_pager_Alert(in *Alert, out *pager.Alert, s conversion.Scope) error {
	err := autoConvert_v1alpha1_Alert_To_pager_Alert(in, out, s)
	if err != nil {
		return err
	}
	if out.ObjectMeta.Annotations == nil {
		return nil
	}
	title := out.ObjectMeta.Annotations[AnnotationCustomTitle]
	out.Spec.Title = title
	return nil
}

func Convert_pager_AlertSpec_To_v1alpha1_AlertSpec(in *pager.AlertSpec, out *AlertSpec, s conversion.Scope) error {
	return autoConvert_pager_AlertSpec_To_v1alpha1_AlertSpec(in, out, s)
}
