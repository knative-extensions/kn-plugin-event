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

package k8s

import (
	"github.com/spf13/pflag"
	knk8s "knative.dev/client/pkg/k8s"
)

// Params contain Kubernetes specific params, that CLI should comply to.
type Params struct {
	Namespace string
	knk8s.Params
}

// Parse will build k8s.Configurator struct which could be used to initiate
// the k8s.Clients.
func (kp *Params) Parse() *Configurator {
	var ns *string
	if kp.Namespace != "" {
		ns = &kp.Namespace
	}
	return &Configurator{
		ClientConfig: kp.GetClientConfig,
		Namespace:    ns,
	}
}

func (kp *Params) SetGlobalFlags(flags *pflag.FlagSet) {
	kp.Params.SetFlags(flags)
}

func (kp *Params) SetCommandFlags(flags *pflag.FlagSet) {
	flags.StringVarP(
		&kp.Namespace, "namespace", "n", "",
		`Specify a namespace of sender job to be created, while event is send
within a cluster. To do that kn-event uses a special Job that is deployed to
cluster in namespace dictated by --namespace. If this option isn't specified
a current context namespace will be used.`,
	)
}
