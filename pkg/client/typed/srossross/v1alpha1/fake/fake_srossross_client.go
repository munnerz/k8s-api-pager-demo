/*

MIT License

Copyright (c) 2017 Sean Ross-Ross

See License in the root of this repo.

*/
package fake

import (
	v1alpha1 "github.com/srossross/k8s-test-controller/pkg/client/typed/srossross/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeSrossrossV1alpha1 struct {
	*testing.Fake
}

func (c *FakeSrossrossV1alpha1) Tests(namespace string) v1alpha1.TestInterface {
	return &FakeTests{c, namespace}
}

func (c *FakeSrossrossV1alpha1) TestRuns(namespace string) v1alpha1.TestRunInterface {
	return &FakeTestRuns{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSrossrossV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
