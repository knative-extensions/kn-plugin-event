package cmd

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/cardil/kn-event/internal"
	"github.com/cardil/kn-event/internal/cli"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestVersionSubCommandWithHuman(t *testing.T) {
	versionSubCommandChecks(t, "human", func(in []byte, out interface{}) (err error) {
		pv := out.(*pluginVersionOutput)
		str := string(in)
		r := regexp.MustCompile("([^ ]+) version: (.+)")
		matches := r.FindStringSubmatch(str)
		pv.Name = matches[1]
		pv.Version = matches[2]
		return nil
	})
}

func TestVersionSubCommandWithJson(t *testing.T) {
	versionSubCommandChecks(t, "json", json.Unmarshal)
}

func TestVersionSubCommandWithYaml(t *testing.T) {
	versionSubCommandChecks(t, "yaml", yaml.Unmarshal)
}

type unmarshalFunc func(in []byte, out interface{}) (err error)

func versionSubCommandChecks(t *testing.T, format string, unmarshal unmarshalFunc) {
	rootCmd.SetArgs([]string{"version", "-o", format})
	buf := bytes.NewBuffer([]byte{})
	rootCmd.SetOut(buf)
	assert.NoError(t, rootCmd.Execute())

	pv := pluginVersionOutput{}
	assert.NoError(t, unmarshal(buf.Bytes(), &pv))
	assert.Equal(t, internal.PluginName, pv.Name)
	assert.Equal(t, internal.Version, pv.Version)
}

func TestPresentAsWithInvalidEnum(t *testing.T) {
	_, err := presentAs(pluginVersionOutput{}, cli.YAML+1)
	assert.Error(t, err)
}
