package sender

import (
	"errors"

	"github.com/cardil/kn-event/internal/event"
	"github.com/cardil/kn-event/internal/k8s"
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
