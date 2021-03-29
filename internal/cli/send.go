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
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Send will send CloudEvent to target.
func (c *App) Send(ce cloudevents.Event, target *TargetArgs, options *OptionsArgs) error {
	t, err := createTarget(target, options.WithLogger())
	if err != nil {
		return err
	}
	sender, err := c.Binding.NewSender(t)
	if err != nil {
		return err
	}
	return sender.Send(ce)
}
