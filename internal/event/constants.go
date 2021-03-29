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

package event

import (
	"fmt"

	"github.com/google/uuid"
	"knative.dev/kn-plugin-event/internal"
)

const (
	// DefaultType holds a default type for a event.
	DefaultType = "dev.knative.cli.plugin.event.generic"
)

// DefaultSource holds a default source of an event.
func DefaultSource() string {
	return fmt.Sprintf("%s/%s", internal.PluginName, internal.Version)
}

// NewID creates a new ID for an event.
func NewID() string {
	return uuid.New().String()
}
