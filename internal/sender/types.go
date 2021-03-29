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

package sender

import (
	"errors"

	"knative.dev/kn-plugin-event/internal/event"
	"knative.dev/kn-plugin-event/internal/k8s"
)

var (
	// ErrUnsupportedTargetType is an error if user pass unsupported event target
	// type. Only supporting: reachable or addressable.
	ErrUnsupportedTargetType = errors.New("unsupported target type")

	// ErrCouldntBeSent is an error that will be return in case event that suppose
	// to be sent, couldn't be, for whatever technical reason.
	ErrCouldntBeSent = errors.New("event couldn't be sent")
)

// CreateKubeClients creates k8s.Clients.
type CreateKubeClients func(props *event.Properties) (k8s.Clients, error)

// CreateJobRunner creates a k8s.JobRunner.
type CreateJobRunner func(kube k8s.Clients) k8s.JobRunner

// CreateAddressResolver creates a k8s.ReferenceAddressResolver.
type CreateAddressResolver func(kube k8s.Clients) k8s.ReferenceAddressResolver

// Binding holds injectable dependencies.
type Binding struct {
	CreateJobRunner
	CreateAddressResolver
	CreateKubeClients
}
