package event

import (
	"errors"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ErrNotYetImplemented is an error for not yet implemented code.
var ErrNotYetImplemented = errors.New("not yet implemented")

// Spec holds specification of event to be created.
type Spec struct {
	Type   string
	ID     string
	Source string
	Fields []FieldSpec
}

// FieldSpec holds a specification of a event's data field.
type FieldSpec struct {
	Path  string
	Value interface{}
}

// TargetType specify a type of a event target.
type TargetType int

const (
	// TargetTypeReachable specify a type of event target that is network
	// reachable, and direct HTTP communication can be performed.
	TargetTypeReachable TargetType = iota

	// TargetTypeAddressable represent a type of event target that is cluster
	// private, and direct communication can't be performed. In this case in
	// cluster sender Job will be created to send the event.
	TargetTypeAddressable
)

// AddressableSpec specify destination of a event to be sent, as well as sender
// namespace that should be used to create a sender Job in.
type AddressableSpec struct {
	duckv1.Destination
	SenderNamespace string
}

// Target is a target to send event to.
type Target struct {
	Type           TargetType
	URLVal         *url.URL
	AddressableVal *AddressableSpec
}

// KnPluginOptions holds options inherited to every Kn plugin.
type KnPluginOptions struct {
	// KnConfig holds kn configuration file (default: ~/.config/kn/config.yaml)
	KnConfig string

	// Kubeconfig holds kubectl configuration file (default: ~/.kube/config)
	Kubeconfig string

	// LogHTTP tells if kn-event plugin should log HTTP requests it makes
	LogHTTP bool
}

// Properties holds a general properties.
type Properties struct {
	KnPluginOptions
	Log *zap.SugaredLogger
}

// Sender will send event to specified target.
type Sender interface {
	// Send will send cloudevents.Event to configured target, or return an error
	// if one occur.
	Send(ce cloudevents.Event) error
}
