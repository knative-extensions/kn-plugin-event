//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cetest "github.com/cloudevents/sdk-go/v2/test"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/icmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/kn-plugin-event/test"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/eventshub"
	eventshubassert "knative.dev/reconciler-test/pkg/eventshub/assert"
	"knative.dev/reconciler-test/pkg/feature"
	"sigs.k8s.io/yaml"
)

const (
	issue228Warn = "child pods are preserved by default when jobs are deleted; " +
		"set propagationPolicy=Background to remove them or set " +
		"propagationPolicy=Orphan to suppress this warning"
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
		log = log.With(json("cmd", cmd))
		log.Info("Running")
		result := icmd.RunCmd(cmd)
		if err := result.Compare(icmd.Expected{
			ExitCode: 0,
			Err:      fmt.Sprintf("Event (ID: %s) have been sent.", ev.ID()),
		}); err != nil {
			handleSendErr(ctx, t, err, ev)
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

// handleSendErr will handle the error from sending event.
func handleSendErr(ctx context.Context, t feature.T, err error, ev cloudevents.Event) {
	// TODO: most of this code should be moved to production CLI, so that in case
	//       of send error, a nice, report is produced.
	//       See: https://github.com/knative-extensions/kn-plugin-event/issues/129
	if err == nil {
		return
	}
	kube := kubeclient.Get(ctx)
	ns := environment.FromContext(ctx).Namespace()
	log := logging.FromContext(ctx)
	jobs := kube.BatchV1().Jobs(ns)
	pods := kube.CoreV1().Pods(ns)
	events := kube.CoreV1().Events(ns)
	jlist, kerr := jobs.List(ctx, metav1.ListOptions{
		LabelSelector: "event-id=" + ev.ID(),
	})
	if kerr != nil {
		log.Error(kerr)
	}
	if len(jlist.Items) != 1 {
		t.Fatal(err)
	}
	jobName := jlist.Items[0].Name
	plist, kerr := pods.List(ctx, metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	})
	if kerr != nil {
		log.Error(kerr)
	}
	podLogs := make([]string, 0, len(plist.Items))
	for _, item := range plist.Items {
		var bytes []byte
		bytes, kerr = pods.GetLogs(item.Name, nil).DoRaw(ctx)
		if kerr != nil {
			log.Error(kerr)
		}
		podLogs = append(podLogs, string(bytes))
	}
	podsYaml, merr := yaml.Marshal(plist.Items)
	if merr != nil {
		log.Error(merr)
	}
	elist, eerr := events.List(ctx, metav1.ListOptions{})
	if eerr != nil {
		log.Error(eerr)
	}
	eventsYaml, eerr := yaml.Marshal(elist.Items)
	if eerr != nil {
		log.Error(eerr)
	}
	t.Fatal(err, "\n\nJob logs (", len(plist.Items), "):\n",
		strings.Join(podLogs, "\n---\n"), "\n\nPods:\n",
		string(podsYaml), "\n\nEvents:\n", string(eventsYaml))
}
