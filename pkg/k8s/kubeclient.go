package k8s

import (
	"errors"
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // see: https://github.com/kubernetes/client-go/issues/242
	"k8s.io/client-go/tools/clientcmd"
	knk8s "knative.dev/client/pkg/k8s"
	eventingv1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/messaging/v1"
	servingv1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// ErrNoKubernetesConnection if we can't connect to Kube API server.
var ErrNoKubernetesConnection = errors.New("no Kubernetes connection")

// NewKubeClients creates Clients.
type NewKubeClients func(configurator *Configurator) (Clients, error)

// Configurator for creating the Kube's clients.
type Configurator struct {
	ClientConfig func() (clientcmd.ClientConfig, error)
	Namespace    *string
}

// NewClients creates kubernetes clients.
func NewClients(cfg *Configurator) (Clients, error) {
	cc, err := loadClientConfig(cfg)
	if err != nil {
		return nil, err
	}
	restcfg, rerr := cc.ClientConfig.ClientConfig()
	if rerr != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoKubernetesConnection, rerr)
	}
	typed, err := kubernetes.NewForConfig(restcfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoKubernetesConnection, err)
	}
	// NOTE: No new error can happen, as all connection related issues are already
	//       checked above.
	dyn, _ := dynamic.NewForConfig(restcfg)
	servingclient, _ := servingv1.NewForConfig(restcfg)
	eventingclient, _ := eventingv1.NewForConfig(restcfg)
	messagingclient, _ := messagingv1.NewForConfig(restcfg)
	return &clients{
		clientConfig: cc,
		typed:        typed,
		dynamic:      dyn,
		serving:      servingclient,
		eventing:     eventingclient,
		messaging:    messagingclient,
	}, nil
}

// Clients holds available Kubernetes clients.
type Clients interface {
	Namespace() string
	Typed() kubernetes.Interface
	Dynamic() dynamic.Interface
	Serving() servingv1.ServingV1Interface
	Eventing() eventingv1.EventingV1Interface
	Messaging() messagingv1.MessagingV1Interface
}

func loadClientConfig(cfg *Configurator) (clientConfig, error) {
	if cfg == nil {
		return clientConfig{}, fmt.Errorf("%w: no config", ErrNoKubernetesConnection)
	}
	ccFn := cfg.ClientConfig
	if ccFn == nil {
		kn := knk8s.Params{}
		ccFn = kn.GetClientConfig
	}
	cc, err := ccFn()
	if err != nil {
		return clientConfig{}, fmt.Errorf("%w: %w", ErrNoKubernetesConnection, err)
	}
	ns, _, err := cc.Namespace()
	if err != nil {
		return clientConfig{}, fmt.Errorf("%w: %w", ErrNoKubernetesConnection, err)
	}
	if cfg.Namespace != nil {
		ns = *cfg.Namespace
	}
	return clientConfig{ClientConfig: cc, namespace: ns}, nil
}

type clientConfig struct {
	clientcmd.ClientConfig
	namespace string
}

type clients struct {
	clientConfig
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
