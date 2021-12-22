//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cetest "github.com/cloudevents/sdk-go/v2/test"
	"gotest.tools/v3/icmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/kn-plugin-event/test"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/eventshub"
	"knative.dev/reconciler-test/pkg/eventshub/assert"
	"knative.dev/reconciler-test/pkg/feature"
	"sigs.k8s.io/yaml"
)

// SendEventToClusterLocal returns a feature.Feature that can be reused in other
// test suites.
func SendEventToClusterLocal() *feature.Feature {
	f := feature.NewFeature()
	sinkName := feature.MakeRandomK8sName("sink")
	ev := cetest.FullEvent()
	ev.SetID(feature.MakeRandomK8sName("test-event"))

	f.Setup("deploy sink", eventshub.Install(sinkName, eventshub.StartReceiver))

	f.Alpha("Event").
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
			"--sender-namespace", ns,
			"--field", fmt.Sprintf("data=%s", ev.Data()),
			"--to", fmt.Sprintf("Service:v1:%s", sinkName),
		}
		cmd := test.ResolveKnEventCommand(t).ToIcmd(args...)
		log.Infof("Running command: %v", cmd)
		result := icmd.RunCmd(cmd)
		err := result.Compare(icmd.Expected{
			ExitCode: 0,
			Out:      fmt.Sprintf("Event (ID: %s) have been sent.", ev.ID()),
		})
		handleSendErr(ctx, t, err, ev)
	}
}

func receiveEvent(ev cloudevents.Event, sinkName string) feature.StepFn {
	return assert.OnStore(sinkName).
		MatchEvent(cetest.HasId(ev.ID())).
		Exact(1)
}

// handleSendErr TODO: most of this code should be moved to production CLI, so
//                     that in case of send error, a nice, report is produced.
//                     See: https://github.com/knative-sandbox/kn-plugin-event/issues/129
func handleSendErr(ctx context.Context, t feature.T, err error, ev cloudevents.Event) {
	if err == nil {
		return
	}
	kube := kubeclient.Get(ctx)
	ns := environment.FromContext(ctx).Namespace()
	log := logging.FromContext(ctx)
	jobs := kube.BatchV1().Jobs(ns)
	pods := kube.CoreV1().Pods(ns)
	jlist, kerr := jobs.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("event-id=%s", ev.ID()),
	})
	if kerr != nil {
		log.Error(kerr)
	}
	jobName := jlist.Items[0].Name
	plist, kerr := pods.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
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
	t.Fatal(err, "\n\nJob logs (", len(plist.Items), "):\n",
		strings.Join(podLogs, "\n---\n"), "\n\nPods:\n",
		string(podsYaml))
}
