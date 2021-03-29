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

package ics

import (
	"errors"

	"knative.dev/kn-plugin-event/internal/event"
)

var (
	// ErrCouldntEncode is returned when problem occur while trying to encode an
	// event.
	ErrCouldntEncode = errors.New("couldn't encode an event")
	// ErrCouldntDecode is returned when problem occur while trying to decode an
	// event.
	ErrCouldntDecode = errors.New("couldn't decode an event")
	// ErrCantConfigureICS is returned when problem occur while trying to
	// configure ICS sender.
	ErrCantConfigureICS = errors.New("can't configure ICS sender")
)

// Args holds a list of args for in-cluster-sender.
type Args struct {
	Sink        string
	CeOverrides string
	Event       string
}

// App holds an ICS app binding.
type App struct {
	event.Binding
}
