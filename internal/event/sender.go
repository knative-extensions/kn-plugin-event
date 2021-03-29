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
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// NewSender will create a sender that can send event to cluster.
func (b Binding) NewSender(target *Target) (Sender, error) {
	sender, err := b.CreateSender(target)
	if err != nil {
		return nil, err
	}
	return &sendLogic{Sender: sender, Properties: target.Properties}, nil
}

type sendLogic struct {
	Sender
	*Properties
}

func (l *sendLogic) Send(ce cloudevents.Event) error {
	err := l.Sender.Send(ce)
	if err == nil {
		l.Log.Infof("Event (ID: %s) have been sent.", ce.ID())
	}
	return err
}
