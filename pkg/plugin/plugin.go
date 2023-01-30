package plugin

import (
	"io"

	"github.com/wavesoftware/go-commandline"
	knplugin "knative.dev/client-pkg/pkg/kn/plugin"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

// init makes sure to register plugin as internal one, after import of
// pkg/plugin, as knative cli plugins are expected to do.
func init() { //nolint:gochecknoinits
	knplugin.InternalPlugins = append(knplugin.InternalPlugins, &plugin{})
}

type plugin struct {
	io.Writer
}

func (p plugin) Name() string {
	return metadata.PluginName
}

func (p plugin) Execute(args []string) error {
	opts := []commandline.Option{
		commandline.WithArgs(args...),
	}
	if p.Writer != nil {
		opts = append(opts, commandline.WithOutput(p.Writer))
	}
	return commandline.New(new(cmd.App)).Execute(opts...) //nolint:wrapcheck
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
