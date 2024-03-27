package k8s

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // see: https://github.com/kubernetes/client-go/issues/242
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	eventingv1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/messaging/v1"
	"knative.dev/kn-plugin-event/pkg/event"
	servingv1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// ErrNoKubernetesConnection if can't connect to Kube API server.
var ErrNoKubernetesConnection = errors.New("no Kubernetes connection")

// CreateKubeClient creates kubernetes.Interface.
func CreateKubeClient(props *event.Properties) (Clients, error) {
	cc, err := loadClientConfig(props)
	if err != nil {
		return nil, err
	}
	restcfg := cc.Config
	typed, err := kubernetes.NewForConfig(restcfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexcpected, err)
	}
	dyn, err := dynamic.NewForConfig(restcfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexcpected, err)
	}
	servingclient, err := servingv1.NewForConfig(restcfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexcpected, err)
	}
	eventingclient, err := eventingv1.NewForConfig(restcfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexcpected, err)
	}
	messagingclient, err := messagingv1.NewForConfig(restcfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexcpected, err)
	}
	return &clients{
		ctx:       context.Background(),
		namespace: cc.namespace,
		typed:     typed,
		dynamic:   dyn,
		serving:   servingclient,
		eventing:  eventingclient,
		messaging: messagingclient,
	}, nil
}

// Clients holds available Kubernetes clients.
type Clients interface {
	Namespace() string
	Typed() kubernetes.Interface
	Dynamic() dynamic.Interface
	Context() context.Context
	Serving() servingv1.ServingV1Interface
	Eventing() eventingv1.EventingV1Interface
	Messaging() messagingv1.MessagingV1Interface
}

func loadClientConfig(props *event.Properties) (clientConfig, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	var configOverrides *clientcmd.ConfigOverrides
	if props.Context != "" && props.Cluster != "" {
		configOverrides = &clientcmd.ConfigOverrides{}
		if props.Context != "" {
			configOverrides.CurrentContext = props.Context
		}
		if props.Cluster != "" {
			configOverrides.Context.Cluster = props.Cluster
		}
	}
	if len(props.Path) > 0 {
		loadingRules.ExplicitPath = props.Path
	}
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	cfg, err := cc.ClientConfig()
	if err != nil {
		return clientConfig{}, fmt.Errorf("%w: %w", ErrNoKubernetesConnection, err)
	}
	ns, _, err := cc.Namespace()
	if err != nil {
		return clientConfig{}, fmt.Errorf("%w: %w", ErrNoKubernetesConnection, err)
	}
	return clientConfig{Config: cfg, namespace: ns}, nil
}

type clientConfig struct {
	*rest.Config
	namespace string
}

type clients struct {
	namespace string
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

func (c *clients) Namespace() string {
	return c.namespace
}
