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
	"knative.dev/kn-plugin-event/internal/configuration"
)

var buildCmd = func() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Builds a CloudEvent and print it to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			options.OutWriter = cmd.OutOrStdout()
			options.ErrWriter = cmd.ErrOrStderr()
			cli := configuration.CreateCli()
			ce, err := cli.CreateWithArgs(eventArgs)
			if err != nil {
				return err
			}
			out, err := cli.PresentWith(ce, options.Output)
			if err != nil {
				return err
			}
			cmd.Println(out)
			return nil
		},
	}
	addBuilderFlags(c)
	return c
}()
