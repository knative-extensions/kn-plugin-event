package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	fakedyna "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	fakekube "k8s.io/client-go/kubernetes/fake"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
)

// FakeClients creates K8s clients from a list of objects using fake packages.
type FakeClients struct {
	testing.TB
	Objects []runtime.Object
	kube    kubernetes.Interface
	dyna    dynamic.Interface
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
		assert.NoError(c.TB, kubescheme.AddToScheme(s))
		c.dyna = fakedyna.NewSimpleDynamicClient(s, c.Objects...)
	}
	return c.dyna
}

func (c *FakeClients) Context() context.Context {
	return c.ctx
}
