package event_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	pkglogging "knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/logging"
)

var errTestError = errors.New("test error")

func TestSendingAnEvent(t *testing.T) {
	testCases := []testCase{
		passingCase(),
		failingSend(),
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()
			ctx = logging.WithTestLogger(ctx, t)
			s, err := tt.binding.NewSender(nil, tt.target)
			if err != nil {
				if !errors.Is(err, tt.want) {
					t.Errorf("want: %#v\n got: %#v", tt.want, err)
				}
				return
			}
			var buf *zaptest.Buffer
			buf, ctx = setupLoggingBuffer(ctx)
			got := s.Send(ctx, tt.ce)
			if !errors.Is(got, tt.want) {
				t.Errorf("want: %#v\n got: %#v", tt.want, got)
			}
			if tt.bufTest != nil {
				tt.bufTest(t, buf)
			}
		})
	}
}

func setupLoggingBuffer(ctx context.Context) (*zaptest.Buffer, context.Context) {
	var buf zaptest.Buffer
	cfg := zap.NewDevelopmentConfig()
	enc := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	log := zap.New(zapcore.NewCore(enc, &buf, cfg.Level))
	ctx = pkglogging.WithLogger(ctx, log.Sugar())
	return &buf, ctx
}

func passingCase() testCase {
	ce := cloudevents.NewEvent("1.0")
	ce.SetID("123456")
	target := &event.Target{
		Reference:   nil,
		RelativeURI: "",
	}
	return testCase{
		bufTest: func(t *testing.T, buf *zaptest.Buffer) {
			t.Helper()

			text := buf.String()
			assert.Check(t, strings.Contains(text, "Event (ID: 123456) have been sent."))
		},
		name: "passing",
		ce:   ce,
		binding: event.Binding{
			CreateSender: stubSenderFactory,
		},
		target: target,
	}
}

func failingSend() testCase {
	return testCase{
		name: "failingSend",
		want: errTestError,
		binding: event.Binding{
			CreateSender: func(_ *k8s.Configurator, _ *event.Target) (event.Sender, error) {
				return nil, errTestError
			},
		},
	}
}

type stubSender struct{}

func (m *stubSender) Send(context.Context, cloudevents.Event) error {
	return nil
}

func stubSenderFactory(*k8s.Configurator, *event.Target) (event.Sender, error) {
	return &stubSender{}, nil
}

type testCase struct {
	name    string
	bufTest func(t *testing.T, buf *zaptest.Buffer)
	target  *event.Target
	ce      cloudevents.Event
	want    error
	binding event.Binding
}
