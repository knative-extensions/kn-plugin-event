package plugin_test

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
	knplugin "knative.dev/client/pkg/kn/plugin"
	"knative.dev/kn-plugin-event/pkg/metadata"
	"knative.dev/kn-plugin-event/pkg/plugin"
	_ "knative.dev/kn-plugin-event/pkg/plugin"
)

func TestPluginRegistersAsInternal(t *testing.T) {
	assert.Check(t, len(knplugin.InternalPlugins) > 0)
}

func TestPluginExecutes(t *testing.T) {
	pl := findPlugin()
	assert.Check(t, pl != nil)
	bytes := plugin.WithCapture(func() {
		err := pl.Execute([]string{"version", "-o", "json"})
		assert.NilError(t, err)
	})
	ver := extactVersionFromJsonOutput(t, bytes)
	assert.Equal(t, ver, metadata.Version)
}

func TestPluginDescription(t *testing.T) {
	pl := findPlugin()
	assert.Check(t, pl != nil)
	desc, err := pl.Description()
	assert.NilError(t, err)
	assert.Equal(t, desc, metadata.PluginDescription)
	assert.DeepEqual(t, pl.CommandParts(), []string{metadata.PluginUse})
	assert.Equal(t, pl.Path(), "")
}

func extactVersionFromJsonOutput(tb testing.TB, bytes []byte) string {
	assert.Check(tb, len(bytes) > 0)
	un := map[string]string{}
	err := json.Unmarshal(bytes, &un)
	assert.NilError(tb, err)
	ver, ok := un["version"]
	assert.Check(tb, ok)
	return ver
}

func findPlugin() knplugin.Plugin {
	for _, pl := range knplugin.InternalPlugins {
		if pl.Name() == metadata.PluginName {
			return pl
		}
	}
	return nil
}
