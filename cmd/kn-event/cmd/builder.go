/*
Copyright 2021 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/spf13/cobra"
	"knative.dev/kn-plugin-event/internal/cli"
	"knative.dev/kn-plugin-event/internal/event"
)

var eventArgs = &cli.EventArgs{}

func addBuilderFlags(c *cobra.Command) {
	c.Flags().StringVarP(
		&eventArgs.Type, "type", "t", event.DefaultType,
		"Specify a type of a CloudEvent",
	)
	c.Flags().StringVarP(
		&eventArgs.ID, "id", "i", event.NewID(),
		"Specify a CloudEvent ID",
	)
	c.Flags().StringVarP(
		&eventArgs.Source, "source", "s", event.DefaultSource(),
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
}
