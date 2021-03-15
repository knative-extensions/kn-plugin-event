package ics_test

import (
	"testing"
	"time"

	"github.com/cardil/kn-event/internal/cli/ics"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
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
