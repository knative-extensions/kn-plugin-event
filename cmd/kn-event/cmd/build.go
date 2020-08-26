package cmd

import (
	"github.com/cardil/kn-event/internal/cli"
	"github.com/cardil/kn-event/internal/event"
	"github.com/spf13/cobra"
)

var (
	eventArgs = &cli.EventArgs{}
	buildCmd  = func() *cobra.Command {
		c := &cobra.Command{
			Use:   "build",
			Short: "Builds a CloudEvent and print it to stdout",
			RunE: func(cmd *cobra.Command, args []string) error {
				e, err := cli.CreateWithArgs(eventArgs)
				if err != nil {
					return err
				}
				out, err := cli.PresentWith(e, Output)
				if err != nil {
					return err
				}
				cmd.Println(out)
				return nil
			},
		}
		c.Flags().StringVarP(
			&eventArgs.Type, "type", "t", event.DefaultType,
			"Specify a type of a CloudEvent",
		)
		c.Flags().StringVarP(
			&eventArgs.ID, "id", "i", event.NewID(),
			"Specify a CloudEvent ID",
		)
		c.Flags().StringVarP(
			&eventArgs.Source, "source", "s", event.DefaultSource,
			"Specify a source of an CloudEvent",
		)
		c.Flags().StringArrayVarP(
			&eventArgs.Fields, "field", "f", make([]string, 0),
			`Specify a field for data of an CloudEvent. Field should be specified as 
jsonpath expression followed by equal sign and then a value. Value will be 
resolved to be used in exact type. Example:
"person.age=18".`,
		)
		c.Flags().StringArrayVar(
			&eventArgs.RawFields, "raw-field", make([]string, 0),
			`Specify a raw field for data of an CloudEvent. Raw field should be 
specified as jsonpath expression followed by equal sign and then a value. The 
value will be used as string. Example: "person.name=John".`,
		)
		return c
	}()
)
