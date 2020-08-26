package cli

import (
	"fmt"
	"testing"
	"time"

	"github.com/cardil/kn-event/internal/event"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestPresentWithHumanOutput(t *testing.T) {
	e := exampleEvent(t)
	out, err := PresentWith(e, HumanReadable)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.knative.cli.plugin.event.generic
  source: %s
  id: 99e4f4f6-08ff-4bff-acf1-47f61ded68c9
  time: 2020-08-24T14:01:12.000601161Z
  datacontenttype: application/json
Data,
  {
    "active": true,
    "person": {
      "email": "ksuszyns@example.org",
      "name": "Chris"
    },
    "ping": 123,
    "ref": "321"
  }`, event.DefaultSource), out)
}

func TestPresentWithJson(t *testing.T) {
	e := exampleEvent(t)
	out, err := PresentWith(e, JSON)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`{
  "data": {
    "active": true,
    "person": {
      "email": "ksuszyns@example.org",
      "name": "Chris"
    },
    "ping": 123,
    "ref": "321"
  },
  "datacontenttype": "application/json",
  "id": "99e4f4f6-08ff-4bff-acf1-47f61ded68c9",
  "source": "%s",
  "specversion": "1.0",
  "time": "2020-08-24T14:01:12.000601161Z",
  "type": "dev.knative.cli.plugin.event.generic"
}`, event.DefaultSource), out)
}

func TestPresentWithYaml(t *testing.T) {
	e := exampleEvent(t)
	out, err := PresentWith(e, YAML)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`data:
  active: true
  person:
    email: ksuszyns@example.org
    name: Chris
  ping: 123
  ref: "321"
datacontenttype: application/json
id: 99e4f4f6-08ff-4bff-acf1-47f61ded68c9
source: %s
specversion: "1.0"
time: "2020-08-24T14:01:12.000601161Z"
type: dev.knative.cli.plugin.event.generic
`, event.DefaultSource), out)
}

func TestCreateWithArgs(t *testing.T) {
	id := event.NewID()
	eventType := "org.example.ping"
	eventSource := "/events/ping"
	args := &EventArgs{
		Type:   eventType,
		ID:     id,
		Source: eventSource,
		Fields: []string{
			"person.name=Chris",
			"person.email=ksuszyns@example.com",
			"ping=123",
			"active=true",
		},
		RawFields: []string{"ref=321"},
	}
	actual, err := CreateWithArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, eventType, actual.Type())
	assert.Equal(t, id, actual.ID())
	assert.Equal(t, eventSource, actual.Source())
	expectedData := map[string]interface{}{
		"person": map[string]interface{}{
			"name": "Chris",
			"email": "ksuszyns@example.com",
		},
		"ping": 123.,
		"ref": "321",
		"active": true,
	}
	actualData, err := event.UnmarshalData(actual.Data())
	assert.NoError(t, err)
	assert.EqualValues(t, expectedData, actualData)
	delta := 1_000_000.
	assert.InDelta(t, time.Now().UnixNano(), actual.Time().UnixNano(), delta)
}

func exampleEvent(t *testing.T) *cloudevents.Event {
	e := event.NewDefault()
	e.SetTime(time.Unix(1598277672, 601161))
	e.SetID("99e4f4f6-08ff-4bff-acf1-47f61ded68c9")
	assert.NoError(t, e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"person": map[string]interface{}{
			"name": "Chris",
			"email": "ksuszyns@example.org",
		},
		"ping": 123,
		"active": true,
		"ref": "321",
	}))
	return e
}
