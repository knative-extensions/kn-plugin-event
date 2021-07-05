package cmd

import (
	"github.com/spf13/cobra"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/event"
)

func addBuilderFlags(eventArgs *cli.EventArgs, c *cobra.Command) {
	fs := c.Flags()
	fs.StringVarP(
		&eventArgs.Type, "type", "t", event.DefaultType,
		"Specify a type of a CloudEvent",
	)
	fs.StringVarP(
		&eventArgs.ID, "id", "i", event.NewID(),
		"Specify a CloudEvent ID",
	)
	fs.StringVarP(
		&eventArgs.Source, "source", "s", event.DefaultSource(),
		"Specify a source of an CloudEvent",
	)
	fs.StringArrayVarP(
		&eventArgs.Fields, "field", "f", make([]string, 0),
		`Specify a field for data of an CloudEvent. Field should be specified as 
jsonpath expression followed by equal sign and then a value. Value will be 
resolved to be used in exact type. Example:
"person.age=18".`,
	)
	fs.StringArrayVar(
		&eventArgs.RawFields, "raw-field", make([]string, 0),
		`Specify a raw field for data of an CloudEvent. Raw field should be 
specified as jsonpath expression followed by equal sign and then a value. The 
value will be used as string. Example: "person.name=John".`,
	)
}
