package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/configuration"
	"knative.dev/kn-plugin-event/pkg/event"
)

var (
	// ErrSendTargetValidationFailed is returned if a send target can't pass a
	// validation.
	ErrSendTargetValidationFailed = errors.New("send target validation failed")

	// ErrCantSendEvent is returned if event can't be sent.
	ErrCantSendEvent = errors.New("can't send event")
)

type sendCommand struct {
	target *cli.TargetArgs
	event  *cli.EventArgs
	*App
}

func (s *sendCommand) command() *cobra.Command {
	c := &cobra.Command{
		Use:   "send",
		Short: "Builds and sends a CloudEvent to recipient",
		RunE:  s.run,
	}
	addBuilderFlags(s.event, c)
	c.Flags().StringVarP(
		&s.target.URL, "to-url", "u", "",
		`Specify an URL to send event to. This option can't be used with 
--to option.`,
	)
	c.Flags().StringVarP(
		&s.target.Addressable, "to", "r", "",
		`Specify an addressable resource to send event to. This argument
takes format kind:apiVersion:name for named resources or
kind:apiVersion:labelKey1=value1,labelKey2=value2 for matching via a
label selector. This option can't be used with --to-url option.`,
	)
	c.Flags().StringVarP(
		&s.target.Namespace, "namespace", "n", "",
		`Specify a namespace of addressable resource defined with --to
option. If this option isn't specified a current context namespace will be used
to find addressable resource. This option can't be used with --to-url option.`,
	)
	c.Flags().StringVar(
		&s.target.SenderNamespace, "sender-namespace", "",
		`Specify a namespace of sender job to be created. While using --to
option, event is send within a cluster. To do that kn-event uses a special Job
that is deployed to cluster in namespace dictated by --sender-namespace. If
this option isn't specified a current context namespace will be used. This
option can't be used with --to-url option.`,
	)
	c.Flags().StringVar(
		&s.target.AddressableURI, "addressable-uri", "",
		`Specify an URI of a target addressable resource. If this option
isn't specified target URL will not be changed. This option can't be used with 
--to-url option.`,
	)
	c.PreRunE = func(cmd *cobra.Command, args []string) error {
		err := cli.ValidateTarget(s.target)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrSendTargetValidationFailed, err)
		}
		return nil
	}
	return c
}

func (s *sendCommand) run(cmd *cobra.Command, _ []string) error {
	c := configuration.CreateCli(cmd)
	ce, err := c.CreateWithArgs(s.event)
	if err != nil {
		return cantBuildEventError(err)
	}
	err = c.Send(*ce, *s.target, &s.Options)
	if err != nil {
		return cantSentEvent(err)
	}
	return nil
}

func cantSentEvent(err error) error {
	if errors.Is(err, event.ErrCantSentEvent) {
		return err
	}
	return fmt.Errorf("%w: %w", event.ErrCantSentEvent, err)
}
