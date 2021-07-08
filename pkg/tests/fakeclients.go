package tests

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	fakedyna "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	fakekube "k8s.io/client-go/kubernetes/fake"
	fakeservingv1 "knative.dev/serving/pkg/client/clientset/versioned/fake"
	servingv1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// FakeClients creates K8s clients from a list of objects using fake packages.
type FakeClients struct {
	Objects []runtime.Object
	kube    kubernetes.Interface
	dyna    dynamic.Interface
	serving servingv1.ServingV1Interface
	ctx     context.Context
}

func (c *FakeClients) Typed() kubernetes.Interface {
	if c.kube == nil {
		c.kube = fakekube.NewSimpleClientset(c.Objects...)
	}
	return c.kube
}

func (c *FakeClients) Dynamic() dynamic.Interface {
	if c.dyna == nil {
		s := runtime.NewScheme()
		c.dyna = fakedyna.NewSimpleDynamicClient(s, c.Objects...)
	}
	return c.dyna
}

func (c *FakeClients) Serving() servingv1.ServingV1Interface {
	if c.serving == nil {
		c.serving = fakeservingv1.NewSimpleClientset(c.Objects...).ServingV1()
	}
	return c.serving
}

func (c *FakeClients) Context() context.Context {
	return c.ctx
}
