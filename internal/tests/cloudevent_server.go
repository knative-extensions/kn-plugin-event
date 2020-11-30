package tests

import (
	"context"
	"net/http/httptest"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// WithCloudEventsServer is a testing utility that help by starting a CloudEvents
// HTTP server which can catch a sent event.
func WithCloudEventsServer(test func(serverURL url.URL) error) (*cloudevents.Event, error) {
	var ce *cloudevents.Event
	receive := func(ctx context.Context, event cloudevents.Event) {
		ce = &event
	}
	ctx := context.Background()
	protocol, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, err
	}
	handler, err := cloudevents.NewHTTPReceiveHandler(ctx, protocol, receive)
	if err != nil {
		return nil, err
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	u, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}
	err = test(*u)
	return ce, err
}
