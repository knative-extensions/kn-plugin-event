package sender_test

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cetest "github.com/cloudevents/sdk-go/v2/test"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"gotest.tools/v3/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/validation"
	"knative.dev/client/pkg/flags/sink"
	"knative.dev/eventing/test/rekt/resources/broker"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/sender"
	"knative.dev/kn-plugin-event/pkg/tests"
	"knative.dev/pkg/logging"
)

const toLongForRFC1123 = 64

var errExampleValidationFault = errors.New("example validation fault")

func TestInClusterSenderSend(t *testing.T) {
	testCases := []inClusterTestCase{
		passingInClusterSenderSend(t),
		couldResolveAddress(t),
		idViolatesRFC1123(t),
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			createKubeClient := func(*k8s.Configurator) (k8s.Clients, error) {
				return &tests.FakeClients{}, nil
			}
			createJobRunner := func(k8s.Clients) k8s.JobRunner {
				return tt.fields.jobRunner
			}
			createAddressResolver := func(k8s.Clients) k8s.ReferenceAddressResolver {
				return tt.fields.addressResolver
			}
			binding := sender.Binding{
				NewKubeClients:     createKubeClient,
				NewJobRunner:       createJobRunner,
				NewAddressResolver: createAddressResolver,
			}
			cfg := &k8s.Configurator{}
			s, err := binding.New(cfg, &event.Target{Reference: tt.fields.reference})
			assert.NilError(t, err)
			log := zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
			ctx := logging.WithLogger(context.TODO(), log.Sugar())
			if err = s.Send(ctx, tt.args.ce); !errors.Is(err, tt.err) {
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
			reference: exampleBrokerReference(t),
			jobRunner: stubJobRunner(func(job *batchv1.Job) bool {
				if snk, ok := envof(job.Spec.Template.Spec.Containers[0].Env, "K_SINK"); ok {
					if snk != "default.demo.brokers.cluster.local" {
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
	sar.isValid = func(*sink.Reference) error {
		return errExampleValidationFault
	}
	return inClusterTestCase{
		name: "couldResolveAddress",
		fields: fields{
			reference:       exampleBrokerReference(t),
			addressResolver: sar,
			jobRunner: stubJobRunner(func(*batchv1.Job) bool {
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
			reference:       exampleBrokerReference(t),
			addressResolver: stubAddressResolver(),
			jobRunner: fnJobRunner(func(_ context.Context, job *batchv1.Job) error {
				name := job.GetName()
				errs := validation.IsDNS1035Label(name)
				if len(errs) > 0 {
					//goland:noinspection GoErrorStringFormat
					return fmt.Errorf("Job.batch \"%s\" is invalid: "+ //nolint:err113
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

type fnJobRunner func(_ context.Context, job *batchv1.Job) error

func (f fnJobRunner) Run(ctx context.Context, job *batchv1.Job) error {
	return f(ctx, job)
}

type ar struct {
	isValid func(ref *sink.Reference) error
}

func (a *ar) ResolveAddress(_ context.Context, ref *sink.Reference, _ string) (*url.URL, error) {
	if a.isValid != nil {
		if err := a.isValid(ref); err != nil {
			return nil, err
		}
	}
	u, err := url.Parse(fmt.Sprintf("%s.%s.%s.cluster.local",
		ref.Name, ref.Namespace, ref.GVR.Resource))
	if err != nil {
		return nil, fmt.Errorf("bad url: %w", err)
	}
	return u, nil
}

func stubJobRunner(isValid func(job *batchv1.Job) bool) k8s.JobRunner {
	return fnJobRunner(func(_ context.Context, job *batchv1.Job) error {
		if !isValid(job) {
			return event.ErrCantSentEvent
		}
		return nil
	})
}

func stubAddressResolver() *ar {
	return &ar{}
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

func exampleBrokerReference(t *testing.T) *sink.Reference {
	t.Helper()
	return &sink.Reference{
		KubeReference: &sink.KubeReference{
			GVR:       broker.GVR(),
			Name:      "default",
			Namespace: "demo",
		},
	}
}

type fields struct {
	reference       *sink.Reference
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
