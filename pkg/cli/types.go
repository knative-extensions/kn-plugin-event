package cli

import (
	"github.com/thediveo/enumflag"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
)

// Params holds a general args for all commands.
type Params struct {
	// OutputMode define the type of output commands should be producing.
	OutputMode

	// Verbose tells should commands display additional information about
	// what's happening? Verbose information is printed on stderr.
	Verbose bool

	// Kubernetes related parameters.
	k8s.Params
}

// EventArgs holds args of event to be created with.
type EventArgs struct {
	Type      string
	ID        string
	Source    string
	Fields    []string
	RawFields []string
}

// TargetArgs holds args specific for even sending.
type TargetArgs struct {
	Sink           string
	AddressableURI string
}

// OutputMode is type of output to produce.
type OutputMode enumflag.Flag

// OutputMode enumeration values.
const (
	HumanReadable OutputMode = iota
	JSON
	YAML
)

// App object.
type App struct {
	event.Binding
}
