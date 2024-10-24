//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"path"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cetest "github.com/cloudevents/sdk-go/v2/test"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/icmd"
	configdir "knative.dev/client/pkg/config/dir"
	"knative.dev/kn-plugin-event/test"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/eventshub"
	eventshubassert "knative.dev/reconciler-test/pkg/eventshub/assert"
	"knative.dev/reconciler-test/pkg/feature"
)

const (
	issue228Warn = "child pods are preserved by default when jobs are deleted; " +
		"set propagationPolicy=Background to remove them or set " +
		"propagationPolicy=Orphan to suppress this warning"
	dirPerm = 0o750
)

// SendEventFeature will create a feature.Feature that will test sending an
// event using in cluster sender to SystemUnderTest.
func SendEventFeature(sut SystemUnderTest) *feature.Feature {
	f := feature.NewFeatureNamed(sut.Name())
	sinkName := feature.MakeRandomK8sName("sink")
	ev := cetest.FullEvent()
	ev.SetID(feature.MakeRandomK8sName("test-event"))

	f.Setup("Deploy EventsHub Sink", eventshub.Install(sinkName, eventshub.StartReceiver))

	sink := sut.Deploy(f, sinkName)

	f.Alpha("Event").
		Must("send", sendEvent(ev, sink)).
		Must("receive", receiveEvent(ev, sinkName))

	return f
}

func sendEvent(ev cloudevents.Event, sink Sink) feature.StepFn {
	return func(ctx context.Context, t feature.T) {
		log := logging.FromContext(ctx).
			With(json("event-id", ev.ID()))

		ns := environment.FromContext(ctx).Namespace()
		args := []string{
			"send",
			"--id", ev.ID(),
			"--source", ev.Source(),
			"--type", ev.Type(),
			"--namespace", ns,
			"--field", fmt.Sprintf("data=%s", ev.Data()),
			"--to", sink.String(),
		}
		cmd := test.ResolveKnEventCommand(t).ToIcmd(args...)
		artifacts := os.Getenv("ARTIFACTS")
		if artifacts == "" {
			artifacts = os.TempDir()
		}
		cacheDir := path.Join(artifacts, t.Name())
		if err := os.MkdirAll(cacheDir, dirPerm); err != nil {
			t.Fatal(err)
		}
		cmd.Env = append(os.Environ(), configdir.CacheDirEnvName+"="+cacheDir)
		log = log.With(json("cmd", cmd))
		log.Info("Running")
		result := icmd.RunCmd(cmd)
		if err := result.Compare(icmd.Expected{
			ExitCode: 0,
			Err:      fmt.Sprintf("Event (ID: %s) have been sent.", ev.ID()),
		}); err != nil {
			t.Fatal(err, "\n\nExecution log: "+
				path.Join(cacheDir, "last-exec.log.jsonl"))
		}
		assert.NotContains(t, result.Stderr(), issue228Warn)
		log.Info("Succeeded")
	}
}

func receiveEvent(ev cloudevents.Event, sinkName string) feature.StepFn {
	return eventshubassert.OnStore(sinkName).
		MatchEvent(cetest.HasId(ev.ID())).
		Exact(1)
}
