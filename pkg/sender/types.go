package sender

import (
	"errors"

	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
)

// ErrUnsupportedTargetType is an error if user pass unsupported event target
// type. Only supporting: reachable or addressable.
var ErrUnsupportedTargetType = errors.New("unsupported target type")

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
