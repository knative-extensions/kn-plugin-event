package k8s

import (
	"context"
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/pkg/signals"
	servingv1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// CreateKubeClient creates kubernetes.Interface.
func CreateKubeClient(props *event.Properties) (Clients, error) {
	config, err := clientcmd.BuildConfigFromFlags("", props.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	typed, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	servingclient, err := servingv1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	return &clients{
		ctx:     signals.NewContext(),
		typed:   typed,
		dynamic: dyn,
		serving: servingclient,
	}, nil
}

// Clients holds available Kubernetes clients.
type Clients interface {
	Typed() kubernetes.Interface
	Dynamic() dynamic.Interface
	Context() context.Context
	Serving() servingv1.ServingV1Interface
}

type clients struct {
	ctx     context.Context
	typed   kubernetes.Interface
	dynamic dynamic.Interface
	serving servingv1.ServingV1Interface
}

func (c *clients) Typed() kubernetes.Interface {
	return c.typed
}

func (c *clients) Dynamic() dynamic.Interface {
	return c.dynamic
}

func (c *clients) Context() context.Context {
	return c.ctx
}

func (c *clients) Serving() servingv1.ServingV1Interface {
	return c.serving
}
