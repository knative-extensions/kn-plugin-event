package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cardil/kn-event/internal"
	"github.com/cardil/kn-event/internal/cli"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type pluginVersionOutput struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the kn event plugin version",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := presentAs(pluginVersionOutput{
			Name:    internal.PluginName,
			Version: internal.Version,
		}, options.Output)
		if err != nil {
			return err
		}
		cmd.Println(output)
		return nil
	},
}

func presentAs(pv pluginVersionOutput, mode cli.OutputMode) (string, error) {
	switch mode {
	case cli.JSON:
		return marshalWith(pv, json.Marshal)
	case cli.YAML:
		return marshalWith(pv, yaml.Marshal)
	case cli.HumanReadable:
		return fmt.Sprintf("%s version: %s", pv.Name, pv.Version), nil
	}
	return "", fmt.Errorf("unsupported mode: %v", mode)
}

type marshalFunc func(in interface{}) (out []byte, err error)

func marshalWith(pv pluginVersionOutput, marchaller marshalFunc) (string, error) {
	bytes, err := marchaller(pv)
	return string(bytes), err
}
