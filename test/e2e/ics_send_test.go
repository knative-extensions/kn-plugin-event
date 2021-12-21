//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"fmt"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cetest "github.com/cloudevents/sdk-go/v2/test"
	"gotest.tools/v3/icmd"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/eventshub"
	"knative.dev/reconciler-test/pkg/eventshub/assert"
	"knative.dev/reconciler-test/pkg/feature"
	reconcilertestk8s "knative.dev/reconciler-test/pkg/k8s"
	"knative.dev/reconciler-test/pkg/knative"
)

func TestInClusterSender(t *testing.T) {
	ctx, env := global.Environment(
		knative.WithKnativeNamespace("knative-eventing"),
		knative.WithLoggingConfig,
		knative.WithTracingConfig,
		reconcilertestk8s.WithEventListener,
		environment.Managed(t),
	)

	env.Test(ctx, t, SendEventToClusterLocal())
}

func SendEventToClusterLocal() *feature.Feature {
	f := feature.NewFeature()
	sinkName := feature.MakeRandomK8sName("sink")
	ev := cetest.FullEvent()

	f.Setup("install sink", eventshub.Install(
		sinkName, eventshub.StartReceiver,
	))

	f.Setup("send event", sendEvent(ev, sinkName))

	f.Assert("receive event on sink", assert.OnStore(sinkName).
		MatchEvent(cetest.HasId(ev.ID())).Exact(1))

	return f
}

func sendEvent(ev cloudevents.Event, sinkName string) feature.StepFn {
	return func(ctx context.Context, t feature.T) {
		cmd := icmd.Command("build/_output/bin/kn-event",
			"send",
			"--id", ev.ID(),
			"--source", ev.Source(),
			"--type", ev.Type(),
			"--field", fmt.Sprintf("msg=%s", ev.Data()),
			"--to", fmt.Sprintf("Service:serving.knative.dev/v1:%s", sinkName),
		)
		result := icmd.RunCmd(cmd)
		result.Assert(t, icmd.Expected{
			ExitCode: 0,
			Out:      "event sent successfully",
		})
	}
}
