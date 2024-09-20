package binding_test

import (
	"os"
	"path"
	"testing"

	"gotest.tools/v3/assert"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	knk8s "knative.dev/client/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/binding"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"sigs.k8s.io/yaml"
)

func TestMemoizeKubeClients(t *testing.T) {
	t.Parallel()
	tcs := []testMemoizeKubeClientsCase{{
		name: "cli",
		fn:   func() event.Binding { return binding.CliApp().Binding },
	}, {
		name: "ics",
		fn:   func() event.Binding { return binding.IcsApp().Binding },
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfgfile := tempConfigFile(t)
			params := k8s.Params{
				Params: knk8s.Params{
					KubeCfgPath: cfgfile,
				},
			}
			cfg := &k8s.Configurator{
				ClientConfig: params.GetClientConfig,
			}
			b := tc.fn()
			cl, err := b.NewKubeClients(cfg)
			assert.NilError(t, err)
			cl2, err2 := b.NewKubeClients(cfg)
			assert.NilError(t, err2)
			assert.Equal(t, cl, cl2)

			ns1 := cl.Namespace()
			assert.Equal(t, "expected", ns1)
			updateConfig(t, cfgfile, func(cfg *clientcmdapi.Config) {
				cfg.Contexts[cfg.CurrentContext].Namespace = "replaced"
			})
			ns2 := cl.Namespace()
			assert.Equal(t, ns1, ns2)
			cl, err = b.NewKubeClients(cfg)
			assert.NilError(t, err)
			ns3 := cl.Namespace()
			assert.Equal(t, ns1, ns3)

			b = tc.fn()
			cl, err = b.NewKubeClients(cfg)
			assert.NilError(t, err)
			ns4 := cl.Namespace()
			assert.Equal(t, "replaced", ns4)
		})
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
	tmpfile := path.Join(tb.TempDir(), "kubeconfig")
	cfg := stubConfig()
	safeConfig(tb, tmpfile, cfg)
	return tmpfile
}

func safeConfig(tb testing.TB, cfgfile string, config clientcmdapi.Config) {
	tb.Helper()
	config.SetGroupVersionKind(clientcmdapi.SchemeGroupVersion.WithKind("Config"))
	bytes, err := yaml.Marshal(config)
	assert.NilError(tb, err)
	err = os.WriteFile(cfgfile, bytes, 0o600)
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
