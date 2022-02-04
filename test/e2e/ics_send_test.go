//go:build e2e
// +build e2e

package e2e_test

import (
	"testing"

	"knative.dev/kn-plugin-event/test"
	"knative.dev/kn-plugin-event/test/e2e"
	"knative.dev/reconciler-test/pkg/environment"
	reconcilertestk8s "knative.dev/reconciler-test/pkg/k8s"
	"knative.dev/reconciler-test/pkg/knative"
)

func TestInClusterSender(t *testing.T) {
	test.MaybeSkip(t)
	e2e.RegisterPackages()

	t.Parallel()

	ctx, env := global.Environment(
		environment.Managed(t),
		reconcilertestk8s.WithEventListener,
		knative.WithKnativeNamespace("knative-eventing"),
		knative.WithLoggingConfig,
		knative.WithTracingConfig,
	)

	env.Test(ctx, t, e2e.SendEventToKubeService())
	env.Test(ctx, t, e2e.SendEventToKnService())
}
