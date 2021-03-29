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

package cli

import (
	"io"

	"github.com/thediveo/enumflag"
	"knative.dev/kn-plugin-event/internal/event"
)

// OptionsArgs holds a general args for all commands.
type OptionsArgs struct {
	event.KnPluginOptions

	// Output define type of output commands should be producing.
	Output OutputMode

	// Verbose tells does commands should display additional information about
	// what's happening? Verbose information is printed on stderr.
	Verbose bool

	OutWriter io.Writer
	ErrWriter io.Writer
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
}
