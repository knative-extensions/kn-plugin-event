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

package event_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/event"
)

func TestCreateWithArgs(t *testing.T) {
	id := event.NewID()
	eventType := "org.example.kn.event.ping"
	eventSource := "/k8s/events/ping"
	spec := &event.Spec{
		Type:   eventType,
		ID:     id,
		Source: eventSource,
		Fields: []event.FieldSpec{
			{Path: "person.name", Value: "Chris"},
			{Path: "person.email", Value: "ksuszyns@example.com"},
			{Path: "ping", Value: 123},
			{Path: "active", Value: true},
			{Path: "ref", Value: "321"},
		},
	}
	actual, err := event.CreateFromSpec(spec)
	assert.NoError(t, err)
	assert.Equal(t, eventType, actual.Type())
	assert.Equal(t, id, actual.ID())
	assert.Equal(t, eventSource, actual.Source())
	expectedData := map[string]interface{}{
		"person": map[string]interface{}{
			"name":  "Chris",
			"email": "ksuszyns@example.com",
		},
		"ping":   123.,
		"ref":    "321",
		"active": true,
	}
	actualData, err := event.UnmarshalData(actual.Data())
	assert.NoError(t, err)
	assert.EqualValues(t, expectedData, actualData)
	delta := 1_000_000.
	assert.InDelta(t, time.Now().UnixNano(), actual.Time().UnixNano(), delta)
}
