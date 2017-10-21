package run

import (
	"fmt"
	"log"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	errors "k8s.io/apimachinery/pkg/util/errors"
	wait "k8s.io/apimachinery/pkg/util/wait"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	Srossross "github.com/srossross/k8s-test-controller/pkg/apis/pager"
)

// TestRunCRDName FIXME: could generate this ?
var TestRunCRDName = "testruns.srossross.github.io"

// TestCRDName FIXME: could generate this ?
var TestCRDName = "tests.srossross.github.io"

// TestRunCRD exposes the testrun as a crd
var TestRunCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: TestRunCRDName,
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   Srossross.GroupName,
		Version: "v1alpha1",
		Scope:   apiextensionsv1beta1.NamespaceScoped,
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:     "testruns",
			Kind:       "TestRun",
			ShortNames: []string{"tr"},
		},
	},
}

// TestCRD exposes a test as a crd
var TestCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: TestCRDName,
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   Srossross.GroupName,
		Version: "v1alpha1",
		Scope:   apiextensionsv1beta1.NamespaceScoped,
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:     "tests",
			Kind:       "Test",
			ShortNames: []string{"tst"},
		},
	},
}

// InstallAllCRDs and wait for them to be ready
func InstallAllCRDs(clientset *apiextensionsclient.Clientset) error {
	var err error

	_, err = InstallCRD(clientset, TestRunCRD)

	if err != nil {
		return err
	}

	_, err = InstallCRD(clientset, TestCRD)

	return err
}

// InstallCRD and wait for it to be ready
func InstallCRD(clientset *apiextensionsclient.Clientset, crdDef *apiextensionsv1beta1.CustomResourceDefinition) (*apiextensionsv1beta1.CustomResourceDefinition, error) {

	log.Printf("Ensure CRD '%v'", crdDef.Name)
	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crdDef)

	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.Printf("       CRD '%v' Already exists", crdDef.Name)
			return crdDef, nil
		}
		return nil, err
	}

	var crd *apiextensionsv1beta1.CustomResourceDefinition

	log.Printf("       Waiting for '%v' to be Established", crdDef.Name)
	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdDef.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {

					return true, err
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					fmt.Printf("Name conflict: %v\n", cond.Reason)
				}
			}
		}
		return false, err
	})
	if err != nil {
		deleteErr := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crdDef.Name, nil)
		if deleteErr != nil {
			return nil, errors.NewAggregate([]error{err, deleteErr})
		}
		return nil, err
	}
	log.Printf("       CRD '%v' Ready", crdDef.Name)
	return crd, nil
}
