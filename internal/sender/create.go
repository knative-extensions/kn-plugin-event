package sender

import (
	"fmt"

	"github.com/cardil/kn-event/internal/event"
)

// New creates a new Sender.
func New(target *event.Target) (event.Sender, error) {
	switch target.Type {
	case event.TargetTypeReachable:
		return &directSender{
			url: *target.URLVal,
		}, nil
	case event.TargetTypeAddressable:
		return &inClusterSender{}, nil
	}
	return nil, fmt.Errorf("%w: %v", ErrUnsupportedTargetType, target.Type)
}
