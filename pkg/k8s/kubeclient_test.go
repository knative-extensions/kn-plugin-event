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

package k8s_test

import (
	"testing"

	"gotest.tools/v3/assert"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/kn-plugin-event/pkg/k8s"
)

func TestNewClients(t *testing.T) {
	basic, err := clientcmd.NewClientConfigFromBytes([]byte(basicKubeconfig))
	assert.NilError(t, err)
	t.Setenv("KUBECONFIG", "/var/blackhole/not-existing.yaml")
	tcs := []newClientsTestCase{{
		name:    "nil",
		wantErr: k8s.ErrNoKubernetesConnection,
	}, {
		name:    "invalid configurator",
		config:  &k8s.Configurator{},
		wantErr: k8s.ErrNoKubernetesConnection,
	}, {
		name: "provided configurator",
		config: &k8s.Configurator{
			ClientConfig: just(basic),
		},
	}}
	for _, tc := range tcs {
		t.Run(tc.name, tc.test)
	}
}

func just(cc clientcmd.ClientConfig) func() (clientcmd.ClientConfig, error) {
	return func() (clientcmd.ClientConfig, error) {
		return cc, nil
	}
}

const basicKubeconfig = `apiVersion: v1
kind: Config
preferences: {}
users:
- name: a
  user:
    client-certificate-data: ""
    client-key-data: ""
clusters:
- name: a
  cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:8080
contexts:
- name: a
  context:
    cluster: a
    user: a
current-context: a
`

type newClientsTestCase struct {
	name    string
	config  *k8s.Configurator
	wantErr error
}

func (tc *newClientsTestCase) test(t *testing.T) {
	cl, err := k8s.NewClients(tc.config)
	if tc.wantErr != nil {
		assert.ErrorIs(t, err, tc.wantErr)
	} else {
		ns := cl.Namespace()
		assert.Check(t, ns != "")
		if tc.config.Namespace != nil {
			assert.Equal(t, *tc.config.Namespace, ns)
		}
		assert.Check(t, cl.Dynamic() != nil)
		assert.Check(t, cl.Eventing() != nil)
		assert.Check(t, cl.Messaging() != nil)
		assert.Check(t, cl.Serving() != nil)
		assert.Check(t, cl.Typed() != nil)
	}
}
