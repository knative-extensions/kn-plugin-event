package configuration_test

import (
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/v3/assert"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"knative.dev/kn-plugin-event/pkg/configuration"
	"knative.dev/kn-plugin-event/pkg/event"
	"sigs.k8s.io/yaml"
)

func TestMemoizeKubeClients(t *testing.T) {
	t.Parallel()
	testMemoizeKubeClientsCases(func(tc testMemoizeKubeClientsCase) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfgfile := tempConfigFile(t)
			defer func() {
				assert.NilError(t, os.Remove(cfgfile))
			}()
			props := &event.Properties{
				KnPluginOptions: event.KnPluginOptions{
					KubeconfigOptions: event.KubeconfigOptions{Path: cfgfile},
				},
			}
			app := tc.fn()
			ns1, err := app.DefaultNamespace(props)
			assert.NilError(t, err)
			updateConfig(t, cfgfile, func(cfg *clientcmdapi.Config) {
				cfg.Contexts[cfg.CurrentContext].Namespace = "replaced"
			})
			ns2, err := app.DefaultNamespace(props)
			assert.NilError(t, err)
			assert.Equal(t, ns1, ns2)
		})
	})
}

func testMemoizeKubeClientsCases(fn func(tc testMemoizeKubeClientsCase)) {
	tcs := []testMemoizeKubeClientsCase{{
		name: "cli",
		fn: func() event.Binding {
			return configuration.CreateCli().Binding
		},
	}, {
		name: "ics",
		fn: func() event.Binding {
			return configuration.CreateIcs().Binding
		},
	}}
	for _, tc := range tcs {
		tc := tc
		fn(tc)
	}
}

type testMemoizeKubeClientsCase struct {
	name string
	fn   func() event.Binding
}

func updateConfig(tb testing.TB, cfgfile string, fn func(cfg *clientcmdapi.Config)) {
	tb.Helper()
	cfg := stubConfig()
	fn(&cfg)
	safeConfig(tb, cfgfile, cfg)
}

func tempConfigFile(tb testing.TB) string {
	tb.Helper()
	tmpfile, err := ioutil.TempFile("", "kubeconfig")
	assert.NilError(tb, err)
	assert.NilError(tb, tmpfile.Close())
	cfg := stubConfig()
	safeConfig(tb, tmpfile.Name(), cfg)
	return tmpfile.Name()
}

func safeConfig(tb testing.TB, cfgfile string, config clientcmdapi.Config) {
	tb.Helper()
	config.SetGroupVersionKind(clientcmdapi.SchemeGroupVersion.WithKind("Config"))
	bytes, err := yaml.Marshal(config)
	assert.NilError(tb, err)
	err = ioutil.WriteFile(cfgfile, bytes, 0600)
	assert.NilError(tb, err)
}

func stubConfig() clientcmdapi.Config {
	return clientcmdapi.Config{
		CurrentContext: "test",
		Clusters: map[string]*clientcmdapi.Cluster{
			"test": {
				Server: "https://api.example.localdomain:6443",
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"test": {
				Cluster:   "test",
				AuthInfo:  "test",
				Namespace: "expected",
			},
		},
	}
}
