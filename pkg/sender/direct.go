package sender

import (
	"context"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type directSender struct {
	url url.URL
}

func (d *directSender) Send(ce cloudevents.Event) error {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		return cantSentEvent(err)
	}

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.TODO(), d.url.String())

	// Send that Event.
	err = c.Send(ctx, ce)
	if !cloudevents.IsACK(err) {
		return cantSentEvent(err)
	}

	return nil
}
