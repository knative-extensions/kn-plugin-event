package cmd

import (
	"github.com/cardil/kn-event/internal/cli"
	"github.com/cardil/kn-event/internal/configuration"
	"github.com/spf13/cobra"
)

var (
	target = &cli.TargetArgs{}

	sendCmd = func() *cobra.Command {
		c := &cobra.Command{
			Use:   "send",
			Short: "Builds and sends a CloudEvent to recipient",
			RunE: func(cmd *cobra.Command, args []string) error {
				cli := configuration.CreateCli()
				ce, err := cli.CreateWithArgs(eventArgs)
				if err != nil {
					return err
				}
				return cli.Send(*ce, target, options)
			},
		}
		addBuilderFlags(c)
		c.Flags().StringVarP(
			&target.URL, "to-url", "u", "",
			`Specify an URL to send event to. This option can't be used with 
--to option.`,
		)
		c.Flags().StringVarP(
			&target.Addressable, "to", "r", "",
			`Specify an addressable resource to send event to. This argument
takes format kind:apiVersion:name for named resources or
kind:apiVersion:labelKey1=value1,labelKey2=value2 for matching via a
label selector. This option can't be used with --to-url option.`,
		)
		c.Flags().StringVarP(
			&target.Namespace, "namespace", "n", "",
			`Specify a namespace of addressable resource defined with --to
option. If this option isn't specified a current context namespace will be used
to find addressable resource. This option can't be used with --to-url option.`,
		)
		c.Flags().StringVar(
			&target.SenderNamespace, "sender-namespace", "",
			`Specify a namespace of sender job to be created. While using --to
option, event is send within a cluster. To do that kn-event uses a special Job
that is deployed to cluster in namespace dictated by --sender-namespace. If
this option isn't specified a current context namespace will be used. This
option can't be used with --to-url option.`,
		)
		c.Flags().StringVar(
			&target.AddressableURI, "addressable-uri", "/",
			`Specify an URI of a target addressable resource. If this option
isn't specified a '/' URI will be used. This option can't be used with 
--to-url option.`,
		)
		c.PreRunE = func(cmd *cobra.Command, args []string) error {
			return cli.ValidateTarget(target)
		}
		return c
	}()
)
