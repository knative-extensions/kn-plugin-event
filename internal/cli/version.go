/*
 Copyright 2024 The Knative Authors

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

package cli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

// ErrUnsupportedOutputMode is returned if user passed a unsupported
// output mode.
var ErrUnsupportedOutputMode = errors.New("unsupported mode")

type versionCommand struct {
	*App
}

func (v *versionCommand) command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the kn event plugin version",
		RunE:  v.run,
	}
}

func (v *versionCommand) run(cmd *cobra.Command, _ []string) error {
	output, err := presentAs(cli.PluginVersionOutput{
		Name:    metadata.PluginName,
		Version: metadata.Version,
		Image:   metadata.ResolveImage(),
	}, v.OutputMode)
	if err != nil {
		return err
	}
	cmd.Println(output)
	return nil
}

func presentAs(pv cli.PluginVersionOutput, mode cli.OutputMode) (string, error) {
	switch mode {
	case cli.JSON:
		return marshalWith(pv, json.Marshal)
	case cli.YAML:
		return marshalWith(pv, yaml.Marshal)
	case cli.HumanReadable:
		return fmt.Sprintf("%s version: %s\nsender image: %s",
			pv.Name, pv.Version, pv.Image), nil
	}
	return "", fmt.Errorf("%w: %v", ErrUnsupportedOutputMode, mode)
}

type marshalFunc func(in interface{}) (out []byte, err error)

func marshalWith(pv cli.PluginVersionOutput, marchaller marshalFunc) (string, error) {
	bytes, err := marchaller(pv)
	if err != nil {
		return "", fmt.Errorf("version %w: %w", ErrCantBePresented, err)
	}
	return string(bytes), err
}
