//go:build e2e
// +build e2e

package e2e_test

import (
	"testing"

	"knative.dev/kn-plugin-event/test/e2e"
	"knative.dev/reconciler-test/pkg/environment"
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

	env.Test(ctx, t, e2e.SendEventToClusterLocal())
}
