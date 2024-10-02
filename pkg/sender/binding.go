package sender

import (
	"errors"
	"fmt"

	"knative.dev/client/pkg/flags/sink"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
)

// ErrUnsupportedTargetType is an error if user pass unsupported event target
// type. Only supporting: reachable or reference.
var ErrUnsupportedTargetType = errors.New("unsupported target type")

// NewJobRunner creates a k8s.JobRunner.
type NewJobRunner func(kube k8s.Clients) k8s.JobRunner

// NewAddressResolver creates a k8s.ReferenceAddressResolver.
type NewAddressResolver func(kube k8s.Clients) k8s.ReferenceAddressResolver

// Binding holds injectable dependencies.
type Binding struct {
	NewJobRunner
	NewAddressResolver
	k8s.NewKubeClients
}

// New creates a new Sender.
func (b *Binding) New(cfg *k8s.Configurator, target *event.Target) (event.Sender, error) {
	switch target.Type() {
	case sink.TypeURL:
		return &directSender{
			url: *target.URL,
		}, nil
	case sink.TypeReference:
		kube, err := b.NewKubeClients(cfg)
		if err != nil {
			return nil, err
		}
		jr := b.NewJobRunner(kube)
		ar := b.NewAddressResolver(kube)
		return &inClusterSender{
			namespace:       kube.Namespace(),
			target:          target,
			jobRunner:       jr,
			addressResolver: ar,
		}, nil
	}
	return nil, fmt.Errorf("%w: %v", ErrUnsupportedTargetType, target.Type())
}

func cantSentEvent(err error) error {
	if errors.Is(err, event.ErrCantSentEvent) {
		return err
	}
	return fmt.Errorf("%w: %w", event.ErrCantSentEvent, err)
}
