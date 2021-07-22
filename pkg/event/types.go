package event

import (
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/tracker"
)

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
	*tracker.Reference
	URI             *apis.URL
	SenderNamespace string
}

// Target is a target to send event to.
type Target struct {
	Type           TargetType
	URLVal         *url.URL
	AddressableVal *AddressableSpec
	*Properties
}

// KubeconfigOptions holds options for Kubernetes Client.
type KubeconfigOptions struct {
	Path    string
	Context string
	Cluster string
}

// KnPluginOptions holds options inherited to every Kn plugin.
type KnPluginOptions struct {
	KubeconfigOptions
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

// CreateSender creates a Sender.
type CreateSender func(target *Target) (Sender, error)

// DefaultNamespace returns a default namespace for connected K8s cluster or
// error is namespace can't be determined.
type DefaultNamespace func(props *Properties) (string, error)

// Binding holds injectable dependencies.
type Binding struct {
	CreateSender
	DefaultNamespace
}
