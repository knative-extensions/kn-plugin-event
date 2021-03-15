package tests

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Sender is a mock sender that will record sent events for testing.
type Sender struct {
	Sent []cloudevents.Event
}

// Send will send event to specified target.
func (m *Sender) Send(ce cloudevents.Event) error {
	m.Sent = append(m.Sent, ce)
	return nil
}
