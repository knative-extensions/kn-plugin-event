/*
Copyright 2021 The Knative Authors
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
	"context"
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/kn-plugin-event/internal/event"
)

// CreateKubeClient creates kubernetes.Interface.
func CreateKubeClient(props *event.Properties) (Clients, error) {
	config, err := clientcmd.BuildConfigFromFlags("", props.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	typed, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnexcpected, err)
	}
	return &clients{
		ctx:     context.TODO(),
		typed:   typed,
		dynamic: dyn,
	}, nil
}

// Clients holds available Kubernetes clients.
type Clients interface {
	Typed() kubernetes.Interface
	Dynamic() dynamic.Interface
	Context() context.Context
}

type clients struct {
	ctx     context.Context
	typed   kubernetes.Interface
	dynamic dynamic.Interface
}

func (c *clients) Typed() kubernetes.Interface {
	return c.typed
}

func (c *clients) Dynamic() dynamic.Interface {
	return c.dynamic
}

func (c *clients) Context() context.Context {
	return c.ctx
}
