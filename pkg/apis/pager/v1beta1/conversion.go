package v1beta1

import (
	"k8s.io/apimachinery/pkg/conversion"

	"github.com/munnerz/k8s-api-pager-demo/pkg/apis/pager"
)

func Convert_pager_AlertSpec_To_v1beta1_AlertSpec(in *pager.AlertSpec, out *AlertSpec, s conversion.Scope) error {
	out.Content = in.Message
	return autoConvert_pager_AlertSpec_To_v1beta1_AlertSpec(in, out, s)
}

func Convert_v1beta1_AlertSpec_To_pager_AlertSpec(in *AlertSpec, out *pager.AlertSpec, s conversion.Scope) error {
	out.Message = in.Content
	return autoConvert_v1beta1_AlertSpec_To_pager_AlertSpec(in, out, s)
}
