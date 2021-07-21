package k8s

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	eventingv1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/messaging/v1"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/pkg/signals"
	servingv1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// ErrNoKubernetesConnection if can't connect to Kube API server.
var ErrNoKubernetesConnection = errors.New("no Kubernetes connection")

// CreateKubeClient creates kubernetes.Interface.
func CreateKubeClient(props *event.Properties) (Clients, error) {
	config, err := createRestConfig(props)
	if err != nil {
		return nil, err
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
	eventingclient, err := eventingv1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	messagingclient, err := messagingv1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	return &clients{
		ctx:       signals.NewContext(),
		typed:     typed,
		dynamic:   dyn,
		serving:   servingclient,
		eventing:  eventingclient,
		messaging: messagingclient,
	}, nil
}

// Clients holds available Kubernetes clients.
type Clients interface {
	Typed() kubernetes.Interface
	Dynamic() dynamic.Interface
	Context() context.Context
	Serving() servingv1.ServingV1Interface
	Eventing() eventingv1.EventingV1Interface
	Messaging() messagingv1.MessagingV1Interface
}

func createRestConfig(props *event.Properties) (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	if len(props.Kubeconfig) > 0 {
		loadingRules.ExplicitPath = props.Kubeconfig
	}
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	cfg, err := cc.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoKubernetesConnection, err)
	}
	return cfg, nil
}

type clients struct {
	ctx       context.Context
	typed     kubernetes.Interface
	dynamic   dynamic.Interface
	serving   servingv1.ServingV1Interface
	eventing  eventingv1.EventingV1Interface
	messaging messagingv1.MessagingV1Interface
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

func (c *clients) Eventing() eventingv1.EventingV1Interface {
	return c.eventing
}

func (c *clients) Messaging() messagingv1.MessagingV1Interface {
	return c.messaging
}
