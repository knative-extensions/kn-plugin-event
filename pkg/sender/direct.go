package sender

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"knative.dev/pkg/apis"
)

type directSender struct {
	url apis.URL
}

func (d *directSender) Send(ctx context.Context, ce cloudevents.Event) error {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		return cantSentEvent(err)
	}

	// Set a target.
	ctx = cloudevents.ContextWithTarget(ctx, d.url.String())

	// Send that Event.
	err = c.Send(ctx, ce)
	if !cloudevents.IsACK(err) {
		return cantSentEvent(err)
	}

	return nil
}
