package k8s

import (
	"context"
	"fmt"

	"github.com/cardil/kn-event/internal/event"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	return &clients{
		ctx:     context.TODO(),
		typed:   typed,
		dynamic: dyn,
	}, nil
}

// Clients holds available Kubernetes clients.
type Clients interface {
	Typed() kubernetes.Interface
	Dynamic() dynamic.Interface
	Context() context.Context
}

type clients struct {
	ctx     context.Context
	typed   kubernetes.Interface
	dynamic dynamic.Interface
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
