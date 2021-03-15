package sender_test

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cardil/kn-event/internal/event"
	"github.com/cardil/kn-event/internal/k8s"
	"github.com/cardil/kn-event/internal/sender"
	"github.com/cardil/kn-event/internal/tests"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/tracker"
)

var errExampleValidationFault = errors.New("example validation fault")

func TestInClusterSenderSend(t *testing.T) {
	testCases := []inClusterTestCase{
		passingInClusterSenderSend(t),
		couldResolveAddress(t),
	}
	for i := range testCases {
		tt := testCases[i]
		t.Run(tt.name, func(t *testing.T) {
			createKubeClient := func(_ *event.Properties) (k8s.Clients, error) {
				return &tests.FakeClients{}, nil
			}
			createJobRunner := func(_ k8s.Clients) k8s.JobRunner {
				return tt.fields.jobRunner
			}
			createAddressResolver := func(_ k8s.Clients) k8s.ReferenceAddressResolver {
				return tt.fields.addressResolver
			}
			binding := sender.Binding{
				CreateKubeClients:     createKubeClient,
				CreateJobRunner:       createJobRunner,
				CreateAddressResolver: createAddressResolver,
			}
			s, err := binding.New(&event.Target{
				Type:           event.TargetTypeAddressable,
				AddressableVal: tt.fields.addressable,
			})
			assert.NoError(t, err)
			if err := s.Send(tt.args.ce); !errors.Is(err, tt.err) {
				t.Errorf("Send() error = %v, wantErr = %v", err, tt.err)
			}
		})
	}
}

func passingInClusterSenderSend(t *testing.T) inClusterTestCase {
	t.Helper()
	return inClusterTestCase{
		name: "passing",
		fields: fields{
			addressable: exampleBrokerAddressableSpec(t),
			jobRunner: stubJobRunner(func(job *batchv1.Job) bool {
				if sink, ok := envof(job.Spec.Template.Spec.Containers[0].Env, "K_SINK"); ok {
					if sink != "default.demo.broker.eventing.dev.cluster.local" {
						return false
					}
				} else {
					return false
				}
				if !strings.Contains(job.Spec.Template.Spec.Containers[0].Image, "kn-event-sender") {
					return false
				}
				return true
			}),
			addressResolver: stubAddressResolver(),
		},
		args: args{
			ce: exampleEvent(t),
		},
		err: nil,
	}
}

func couldResolveAddress(t *testing.T) inClusterTestCase {
	t.Helper()
	ar := stubAddressResolver()
	ar.isValid = func(ref *tracker.Reference) error {
		return errExampleValidationFault
	}
	return inClusterTestCase{
		name: "couldResolveAddress",
		fields: fields{
			addressable:     exampleBrokerAddressableSpec(t),
			addressResolver: ar,
			jobRunner: stubJobRunner(func(job *batchv1.Job) bool {
				return true
			}),
		},
		args: args{
			ce: exampleEvent(t),
		},
		err: k8s.ErrInvalidReference,
	}
}

func envof(envs []corev1.EnvVar, name string) (string, bool) {
	for _, env := range envs {
		if env.Name == name {
			return env.Value, true
		}
	}
	return "", false
}

type jr struct {
	isValid func(job *batchv1.Job) bool
}

func (j *jr) Run(job *batchv1.Job) error {
	if !j.isValid(job) {
		return sender.ErrCouldntBeSent
	}
	return nil
}

type ar struct {
	isValid func(ref *tracker.Reference) error
}

func (a *ar) ResolveAddress(ref *tracker.Reference, uri *apis.URL) (*url.URL, error) {
	if a.isValid != nil {
		if err := a.isValid(ref); err != nil {
			return nil, err
		}
	}
	u, err := url.Parse(fmt.Sprintf("%s.%s.%s.cluster.local",
		ref.Name, ref.Namespace, ref.Kind))
	if err != nil {
		return nil, fmt.Errorf("bad url: %w", err)
	}
	return u, nil
}

func stubJobRunner(isValid func(job *batchv1.Job) bool) *jr {
	return &jr{isValid: isValid}
}

func stubAddressResolver() *ar {
	return &ar{}
}

func uri(t *testing.T, uri string) *apis.URL {
	t.Helper()
	u, err := apis.ParseURL(uri)
	assert.NoError(t, err)
	return u
}

func exampleEvent(t *testing.T) cloudevents.Event {
	t.Helper()
	e := cloudevents.NewEvent()
	e.SetID("qazw-sxed-c123")
	e.SetType("testing")
	e.SetSource("source")
	e.SetTime(time.Unix(1_615_678_145, 0))
	assert.NoError(t, e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"person": map[string]interface{}{
			"name":  "Chris",
			"email": "ksuszyns@example.com",
		},
		"ping":   123,
		"active": true,
		"ref":    "321",
	}))
	return e
}

func exampleBrokerAddressableSpec(t *testing.T) *event.AddressableSpec {
	t.Helper()
	return &event.AddressableSpec{
		Reference: &tracker.Reference{
			APIVersion: "betav1",
			Kind:       "broker.eventing.dev",
			Namespace:  "demo",
			Name:       "default",
			Selector:   nil,
		},
		URI:             uri(t, "/"),
		SenderNamespace: "default",
	}
}

type fields struct {
	addressable     *event.AddressableSpec
	addressResolver k8s.ReferenceAddressResolver
	jobRunner       k8s.JobRunner
}

type args struct {
	ce cloudevents.Event
}

type inClusterTestCase struct {
	name   string
	fields fields
	args   args
	err    error
}
