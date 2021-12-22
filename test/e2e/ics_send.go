//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cetest "github.com/cloudevents/sdk-go/v2/test"
	"gotest.tools/v3/icmd"
	"knative.dev/kn-plugin-event/test"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/eventshub"
	"knative.dev/reconciler-test/pkg/eventshub/assert"
	"knative.dev/reconciler-test/pkg/feature"
)

const AwaitForSinkTime = 5 * time.Second

// SendEventToClusterLocal returns a feature.Feature that can be reused in other
// test suites.
func SendEventToClusterLocal() *feature.Feature {
	f := feature.NewFeature()
	sinkName := feature.MakeRandomK8sName("sink")
	ev := cetest.FullEvent()
	ev.SetID(feature.MakeRandomK8sName("test-event"))

	f.Setup("deploy sink", eventshub.Install(sinkName, eventshub.StartReceiver))

	f.Setup("await for sink", func(ctx context.Context, t feature.T) {
		// FIXME: remove the static wait
		time.Sleep(AwaitForSinkTime)
	})

	f.Stable("Event").
		Must("send", sendEvent(ev, sinkName)).
		Must("receive", receiveEvent(ev, sinkName))

	return f
}

func sendEvent(ev cloudevents.Event, sinkName string) feature.StepFn {
	return func(ctx context.Context, t feature.T) {
		log := logging.FromContext(ctx)
		ns := environment.FromContext(ctx).Namespace()
		args := []string{
			"send",
			"--id", ev.ID(),
			"--source", ev.Source(),
			"--type", ev.Type(),
			"--namespace", ns,
			"--field", fmt.Sprintf("data=%s", ev.Data()),
			"--to", fmt.Sprintf("Service:v1:%s", sinkName),
		}
		cmd := test.ResolveKnEventCommand(t).ToIcmd(args...)
		log.Infof("Running command: %v", cmd)
		result := icmd.RunCmd(cmd)
		result.Assert(t, icmd.Expected{
			ExitCode: 0,
			Out:      fmt.Sprintf("Event (ID: %s) have been sent.", ev.ID()),
		})
	}
}

func receiveEvent(ev cloudevents.Event, sinkName string) feature.StepFn {
	return assert.OnStore(sinkName).
		MatchEvent(cetest.HasId(ev.ID())).
		Exact(1)
}
