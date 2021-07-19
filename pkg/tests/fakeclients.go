package tests

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	fakekube "k8s.io/client-go/kubernetes/fake"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	eventingv1fakeclient "knative.dev/eventing/pkg/client/clientset/versioned/fake"
	messagingv1fakeclient "knative.dev/eventing/pkg/client/clientset/versioned/fake"
	eventingv1client "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1"
	messagingv1client "knative.dev/eventing/pkg/client/clientset/versioned/typed/messaging/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingv1fakeclient "knative.dev/serving/pkg/client/clientset/versioned/fake"
	servingv1client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// FakeClients creates K8s clients from a list of objects using fake packages.
type FakeClients struct {
	testing.TB
	Objects   []runtime.Object
	kube      kubernetes.Interface
	dyna      dynamic.Interface
	serving   servingv1client.ServingV1Interface
	eventing  eventingv1client.EventingV1Interface
	messaging messagingv1client.MessagingV1Interface
	ctx       context.Context
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
		assert.NilError(c.TB, kubescheme.AddToScheme(s))
		assert.NilError(c.TB, servingv1.AddToScheme(s))
		assert.NilError(c.TB, eventingv1.AddToScheme(s))
		assert.NilError(c.TB, messagingv1.AddToScheme(s))
		c.dyna = fakedynamic.NewSimpleDynamicClient(s, c.Objects...)
	}
	return c.dyna
}

func (c *FakeClients) Serving() servingv1client.ServingV1Interface {
	if c.serving == nil {
		c.serving = servingv1fakeclient.NewSimpleClientset(c.Objects...).ServingV1()
	}
	return c.serving
}

func (c *FakeClients) Eventing() eventingv1client.EventingV1Interface {
	if c.eventing == nil {
		c.eventing = eventingv1fakeclient.NewSimpleClientset(c.Objects...).EventingV1()
	}
	return c.eventing
}

func (c *FakeClients) Messaging() messagingv1client.MessagingV1Interface {
	if c.messaging == nil {
		c.messaging = messagingv1fakeclient.NewSimpleClientset(c.Objects...).MessagingV1()
	}
	return c.messaging
}

func (c *FakeClients) Context() context.Context {
	if c.ctx == nil {
		c.ctx = context.Background()
	}
	return c.ctx
}
