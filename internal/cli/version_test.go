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

package cli_test

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"gopkg.in/yaml.v2"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

func TestVersionSubCommandWithHuman(t *testing.T) {
	versionSubCommandChecks(t, "human", func(in []byte, out interface{}) error {
		pv, ok := out.(*cli.PluginVersionOutput)
		assert.Check(t, ok)
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

type unmarshalFunc func(in []byte, out interface{}) error

func versionSubCommandChecks(t *testing.T, format string, unmarshal unmarshalFunc) {
	t.Helper()
	buf := bytes.NewBuffer([]byte{})
	assert.NilError(t, testapp().Execute(
		commandline.WithCommand(func(cmd *cobra.Command) {
			cmd.SetOut(buf)
			cmd.SetArgs([]string{"version", "-o", format})
		}),
	))

	pv := cli.PluginVersionOutput{}
	assert.NilError(t, unmarshal(buf.Bytes(), &pv))
	assert.Equal(t, metadata.PluginName, pv.Name)
	assert.Equal(t, metadata.Version, pv.Version)
}

func TestPresentAsWithInvalidOutput(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	err := testapp().Execute(
		commandline.WithCommand(func(cmd *cobra.Command) {
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"version", "-o", "invalid"})
		}),
	)
	assert.Error(t, err, "invalid argument \"invalid\" for "+
		"\"-o, --output\" flag: must be 'human', 'json', 'yaml'")
}
