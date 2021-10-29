package ics_test

import (
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/cli/ics"
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
	assert.NilError(t, err)
	want := &ce

	repr, err := ics.Encode(ce)
	assert.NilError(t, err)
	assert.Check(t, repr != "")
	got, err := ics.Decode(repr)
	assert.NilError(t, err)

	assert.DeepEqual(t, want, got)
}
