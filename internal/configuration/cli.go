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

package configuration

import (
	"knative.dev/kn-plugin-event/internal/cli"
)

// CreateCli creates the configured cli.App to work with.
func CreateCli() *cli.App {
	binding := senderBinding()
	return &cli.App{Binding: eventsBinding(binding)}
}
