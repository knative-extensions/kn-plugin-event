package sender

import (
	"context"
	"fmt"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/wavesoftware/go-ensure"
)

type directSender struct {
	url url.URL
}

func (d *directSender) Send(ce cloudevents.Event) error {
	c, err := cloudevents.NewDefaultClient()
	ensure.NoError(err)

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.TODO(), d.url.String())

	// Send that Event.
	err = c.Send(ctx, ce)
	if !cloudevents.IsACK(err) {
		return fmt.Errorf("%v: %w", ErrCouldntBeSent, err)
	}

	return nil
}
