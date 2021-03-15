package cli

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Send will send CloudEvent to target.
func (c *App) Send(ce cloudevents.Event, target *TargetArgs, options *OptionsArgs) error {
	t, err := createTarget(target, options.WithLogger())
	if err != nil {
		return err
	}
	sender, err := c.Binding.NewSender(t)
	if err != nil {
		return err
	}
	return sender.Send(ce)
}
