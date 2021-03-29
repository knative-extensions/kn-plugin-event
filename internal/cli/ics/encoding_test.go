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

package ics_test

import (
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/cli/ics"
)

func TestEncodeDecode(t *testing.T) {
	ce := cloudevents.NewEvent()
	ce.SetID("987-654-321")
	ce.SetSource("testing://encode-decode")
	ce.SetType("simple")
	ce.SetTime(time.Now().UTC())
	err := ce.SetData("application/json", map[string]interface{}{
		"value": 42,
	})
	assert.NoError(t, err)
	want := &ce

	repr, err := ics.Encode(ce)
	assert.NoError(t, err)
	assert.NotEmpty(t, repr)
	got, err := ics.Decode(repr)
	assert.NoError(t, err)

	assert.EqualValues(t, want, got)
}
