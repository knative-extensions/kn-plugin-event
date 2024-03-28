package sender

import (
	"errors"
	"fmt"

	"knative.dev/kn-plugin-event/pkg/event"
)

// New creates a new Sender.
func (b *Binding) New(target *event.Target) (event.Sender, error) {
	switch target.Type {
	case event.TargetTypeReachable:
		return &directSender{
			url: *target.URLVal,
		}, nil
	case event.TargetTypeAddressable:
		kube, err := b.CreateKubeClients(target.Properties)
		if err != nil {
			return nil, err
		}
		jr := b.CreateJobRunner(kube)
		ar := b.CreateAddressResolver(kube)
		return &inClusterSender{
			addressable:     target.AddressableVal,
			jobRunner:       jr,
			addressResolver: ar,
		}, nil
	}
	return nil, fmt.Errorf("%w: %v", ErrUnsupportedTargetType, target.Type)
}

func cantSentEvent(err error) error {
	if errors.Is(err, event.ErrCantSentEvent) {
		return err
	}
	return fmt.Errorf("%w: %w", event.ErrCantSentEvent, err)
}
