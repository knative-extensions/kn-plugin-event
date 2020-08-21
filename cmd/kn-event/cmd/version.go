package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cardil/kn-event/internal"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
)

type pluginVersionOutput struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the kn event plugin version",
		RunE: func(cmd *cobra.Command, args []string) error {
			pv := pluginVersionOutput{
				Name:    internal.PluginName,
				Version: internal.Version,
			}
			switch Output {
			case HumanReadable:
				cmd.Printf("%s version: %s\n", pv.Name, pv.Version)
			case JSON, YAML:
				bytes, err := marshalWith(pv, Output)
				if err != nil {
					return err
				}
				cmd.Println(string(bytes))
			}
			return nil
		},
	}
)

func marshalWith(in interface{}, mode OutputMode) ([]byte, error) {
	switch mode {
	case JSON:
		return json.Marshal(in)
	case YAML:
		return yaml.Marshal(in)
	}
	return nil, fmt.Errorf("unsupported mode: %v", mode)
}
