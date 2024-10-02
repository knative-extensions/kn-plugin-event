package sender_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/phayes/freeport"
	"knative.dev/client/pkg/flags/sink"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/sender"
	"knative.dev/kn-plugin-event/pkg/tests"
	"knative.dev/pkg/apis"
)

func TestDirectSenderSend(t *testing.T) {
	testsCases := []testCase{
		passing(t),
		undelivered(t),
	}
	for i := range testsCases {
		tt := testsCases[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.context(func(u apis.URL) {
				binding := sender.Binding{}
				cfg := &k8s.Configurator{}
				s, err := binding.New(cfg, &event.Target{Reference: &sink.Reference{
					URL: &u,
				}})
				if err != nil {
					t.Error(err)
					return
				}
				validateErr := func(err error) {
					unexpectedError(t, err)
				}
				if tt.validateErr != nil {
					validateErr = tt.validateErr
				}
				ctx := context.TODO()
				validateErr(s.Send(ctx, tt.ce))
			})
		})
	}
}

func passing(t *testing.T) testCase {
	t.Helper()
	ce := newEvent("543562")
	return testCase{
		name:    "passing",
		ce:      ce,
		context: sentEventIsValid(t, ce),
	}
}

func undelivered(t *testing.T) testCase {
	t.Helper()
	ce := newEvent("1294756")
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Error(err)
		return testCase{}
	}
	u, err := apis.ParseURL(fmt.Sprintf("http://localhost:%d/ce-not-supported", port))
	if err != nil {
		t.Error(err)
		return testCase{}
	}
	return testCase{
		name: "undelivered",
		ce:   ce,
		context: func(handler func(u apis.URL)) {
			handler(*u)
		},
		validateErr: func(err error) {
			if !errors.Is(err, event.ErrCantSentEvent) {
				unexpectedError(t, err)
			}
		},
	}
}

func unexpectedError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error("Send(): unexpected error: ", err)
	}
}

func newEvent(id string) cloudevents.Event {
	ce := cloudevents.NewEvent("1.0")
	ce.SetType("org.example.type")
	ce.SetSource("events://void")
	ce.SetTime(time.Now())
	ce.SetID(id)
	return ce
}

func sentEventIsValid(t *testing.T, want cloudevents.Event) func(hand func(u apis.URL)) {
	t.Helper()
	return func(hand func(u apis.URL)) {
		sent, err := tests.WithCloudEventsServer(func(serverURL url.URL) error {
			hand(apis.URL(serverURL))
			return nil
		})
		if err != nil {
			t.Error(err)
		}
		compareByJSON(t, want, *sent)
	}
}

func compareByJSON(t *testing.T, want interface{}, got interface{}) {
	t.Helper()
	prefix := ""
	indent := "  "
	wantJSON, err := json.MarshalIndent(want, prefix, indent)
	if err != nil {
		t.Error(err)
	}
	gotJSON, err := json.MarshalIndent(got, prefix, indent)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(wantJSON, gotJSON) {
		t.Errorf("events differ:\nwant = %#v\n got = %#v",
			string(wantJSON), string(gotJSON),
		)
	}
}

type testCase struct {
	name        string
	ce          cloudevents.Event
	validateErr func(error)
	context     func(func(u apis.URL))
}
