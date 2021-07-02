package cmd_test

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"knative.dev/kn-plugin-event/cmd/kn-event/cmd"
	"knative.dev/kn-plugin-event/internal"
	"knative.dev/kn-plugin-event/internal/cli"
)

func TestVersionSubCommandWithHuman(t *testing.T) {
	versionSubCommandChecks(t, "human", func(in []byte, out interface{}) (err error) {
		pv, ok := out.(*cli.PluginVersionOutput)
		assert.True(t, ok)
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
	t.Helper()
	tc := cmd.TestingCmd{}
	tc.Args("version", "-o", format)
	buf := bytes.NewBuffer([]byte{})
	tc.Out(buf)
	assert.NoError(t, tc.Execute())

	pv := cli.PluginVersionOutput{}
	assert.NoError(t, unmarshal(buf.Bytes(), &pv))
	assert.Equal(t, internal.PluginName, pv.Name)
	assert.Equal(t, internal.Version, pv.Version)
}

func TestPresentAsWithInvalidOutput(t *testing.T) {
	tc := cmd.TestingCmd{}
	buf := bytes.NewBuffer([]byte{})
	tc.Out(buf)
	tc.Args("version", "-o", "invalid")
	err := tc.Execute()
	assert.IsType(t, cmd.ErrUnsupportedOutputMode, err)
}
