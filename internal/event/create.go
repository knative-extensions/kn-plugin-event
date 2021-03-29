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
	"encoding/json"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/wavesoftware/go-ensure"
)

// NewDefault creates a default CloudEvent.
func NewDefault() *cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetType(DefaultType)
	e.SetID(NewID())
	ensure.NoError(e.SetData(cloudevents.ApplicationJSON, map[string]string{}))
	e.SetSource(DefaultSource())
	e.SetTime(time.Now())
	ensure.NoError(e.Validate())
	return &e
}

// CreateFromSpec will create an event by parsing given args.
func CreateFromSpec(spec *Spec) (*cloudevents.Event, error) {
	e := NewDefault()
	e.SetID(spec.ID)
	e.SetSource(spec.Source)
	e.SetType(spec.Type)
	m := map[string]interface{}{}
	for _, fieldSpec := range spec.Fields {
		updateMapWithSpec(m, fieldSpec)
	}
	err := e.SetData(cloudevents.ApplicationJSON, m)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func updateMapWithSpec(m map[string]interface{}, spec FieldSpec) {
	paths := strings.Split(spec.Path, ".")
	curr := m
	for i, p := range paths {
		if i < len(paths)-1 {
			if _, ok := curr[p]; !ok {
				curr[p] = map[string]interface{}{}
			}
			curr = curr[p].(map[string]interface{})
		} else {
			curr[p] = spec.Value
		}
	}
}

// UnmarshalData will take bytes and unmarshall it as JSON to map structure.
func UnmarshalData(bytes []byte) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// AddField will add a field to the spec.
func (s *Spec) AddField(path string, val interface{}) {
	s.Fields = append(s.Fields, FieldSpec{
		Path: path, Value: val,
	})
}
