package sender

import (
	"github.com/cardil/kn-event/internal/event"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type inClusterSender struct {
}

func (i *inClusterSender) Send(ce cloudevents.Event) error {
	return event.ErrNotYetImplemented
}
