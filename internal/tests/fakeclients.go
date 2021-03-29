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

package tests

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	fakedyna "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	fakekube "k8s.io/client-go/kubernetes/fake"
)

// FakeClients creates K8s clients from a list of objects using fake packages.
type FakeClients struct {
	Objects []runtime.Object
	kube    kubernetes.Interface
	dyna    dynamic.Interface
	ctx     context.Context
}

func (c *FakeClients) Typed() kubernetes.Interface {
	if c.kube == nil {
		c.kube = fakekube.NewSimpleClientset(c.Objects...)
	}
	return c.kube
}

func (c *FakeClients) Dynamic() dynamic.Interface {
	if c.dyna == nil {
		s := runtime.NewScheme()
		c.dyna = fakedyna.NewSimpleDynamicClient(s, c.Objects...)
	}
	return c.dyna
}

func (c *FakeClients) Context() context.Context {
	return c.ctx
}
