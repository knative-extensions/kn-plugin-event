package cli

import (
	"github.com/thediveo/enumflag"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/system"
)

// Options holds a general args for all commands.
type Options struct {
	event.KnPluginOptions

	// Output define type of output commands should be producing.
	Output OutputMode

	// Verbose tells does commands should display additional information about
	// what's happening? Verbose information is printed on stderr.
	Verbose bool
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
	URL             string
	Addressable     string
	Namespace       string
	SenderNamespace string
	AddressableURI  string
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
	system.Environment
}
