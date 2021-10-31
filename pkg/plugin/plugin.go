package plugin

import (
	"io"

	knplugin "knative.dev/client/pkg/kn/plugin"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

// init makes sure to register plugin as internal one, after import of
// pkg/plugin, as knative cli plugins are expected to do.
// nolint:gochecknoinits
func init() {
	knplugin.InternalPlugins = append(knplugin.InternalPlugins, &plugin{})
}

type plugin struct {
	io.Writer
}

func (p plugin) Name() string {
	return metadata.PluginName
}

func (p plugin) Execute(args []string) error {
	opts := []cmd.CommandOption{
		cmd.WithArgs(args...),
	}
	if p.Writer != nil {
		opts = append(opts, cmd.WithOutput(p.Writer))
	}
	return new(cmd.Cmd).ExecuteWithOptions(opts...) //nolint:wrapcheck
}

func (p plugin) Description() (string, error) {
	return metadata.PluginDescription, nil
}

func (p plugin) CommandParts() []string {
	return []string{metadata.PluginUse}
}

func (p plugin) Path() string {
	return ""
}
