package sender_test

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cetest "github.com/cloudevents/sdk-go/v2/test"
	"gotest.tools/v3/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/validation"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/sender"
	"knative.dev/kn-plugin-event/pkg/tests"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/tracker"
)

const toLongForRFC1123 = 64

var errExampleValidationFault = errors.New("example validation fault")

func TestInClusterSenderSend(t *testing.T) {
	testCases := []inClusterTestCase{
		passingInClusterSenderSend(t),
		couldResolveAddress(t),
		idViolatesRFC1123(t),
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
			assert.NilError(t, err)
			if err = s.Send(tt.args.ce); !errors.Is(err, tt.err) {
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
	sar := stubAddressResolver()
	sar.isValid = func(ref *tracker.Reference) error {
		return errExampleValidationFault
	}
	return inClusterTestCase{
		name: "couldResolveAddress",
		fields: fields{
			addressable:     exampleBrokerAddressableSpec(t),
			addressResolver: sar,
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

func idViolatesRFC1123(t *testing.T) inClusterTestCase {
	t.Helper()
	ce := cetest.FullEvent()
	ce.SetID(newIDViolatesRFC1123())
	return inClusterTestCase{
		name: "idViolatesRFC1123",
		fields: fields{
			addressable:     exampleBrokerAddressableSpec(t),
			addressResolver: stubAddressResolver(),
			jobRunner: fnJobRunner(func(job *batchv1.Job) error {
				name := job.GetName()
				errs := validation.IsDNS1035Label(name)
				if len(errs) > 0 {
					//goland:noinspection GoErrorStringFormat
					return fmt.Errorf("Job.batch \"%s\" is invalid: "+ //nolint:goerr113
						"metadata.name: Invalid value: \"%s\": %s",
						name, name, strings.Join(errs, ", "))
				}
				return nil
			}),
		},
		args: args{
			ce: ce,
		},
	}
}

// newIDViolatesRFC1123 returns a new random ID which violates the RFC 1123 on purpose.
func newIDViolatesRFC1123() string {
	return "test-event-" + strings.ToUpper(rand.String(toLongForRFC1123))
}

func envof(envs []corev1.EnvVar, name string) (string, bool) {
	for _, env := range envs {
		if env.Name == name {
			return env.Value, true
		}
	}
	return "", false
}

type fnJobRunner func(job *batchv1.Job) error

func (f fnJobRunner) Run(job *batchv1.Job) error {
	return f(job)
}

type ar struct {
	isValid func(ref *tracker.Reference) error
}

func (a *ar) ResolveAddress(ref *tracker.Reference, _ *apis.URL) (*url.URL, error) {
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

func stubJobRunner(isValid func(job *batchv1.Job) bool) k8s.JobRunner {
	return fnJobRunner(func(job *batchv1.Job) error {
		if !isValid(job) {
			return event.ErrCantSentEvent
		}
		return nil
	})
}

func stubAddressResolver() *ar {
	return &ar{}
}

func uri(t *testing.T, uri string) *apis.URL {
	t.Helper()
	u, err := apis.ParseURL(uri)
	assert.NilError(t, err)
	return u
}

func exampleEvent(t *testing.T) cloudevents.Event {
	t.Helper()
	e := cloudevents.NewEvent()
	e.SetID("qazw-sxed-c123")
	e.SetType("testing")
	e.SetSource("source")
	e.SetTime(time.Unix(1_615_678_145, 0))
	assert.NilError(t, e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
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
