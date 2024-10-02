package k8s_test

import (
	"context"
	"testing"

	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	clienttest "knative.dev/client/pkg/util/test"
	"knative.dev/kn-plugin-event/pkg/k8s"
	k8stest "knative.dev/kn-plugin-event/pkg/k8s/test"
	"knative.dev/kn-plugin-event/pkg/tests"
	"knative.dev/pkg/logging"
)

func TestResolveAddress(t *testing.T) {
	ns := clienttest.NextNamespace()
	for _, tc := range k8stest.ResolveAddressTestCases(ns) {
		t.Run(tc.Name, func(t *testing.T) {
			log := zaptest.NewLogger(t, zaptest.Level(zapcore.WarnLevel))
			ctx := logging.WithLogger(context.TODO(), log.Sugar())
			k8stest.EnsureResolveAddress(ctx, t, tc, func() (k8s.Clients, func(tb testing.TB)) {
				return fakeClients(t, tc), noCleanup
			})
		})
	}
}

func noCleanup(tb testing.TB) {
	tb.Helper() // do nothing
}

func fakeClients(tb testing.TB, tc k8stest.ResolveAddressTestCase) k8s.Clients {
	tb.Helper()
	return &tests.FakeClients{
		Objects: tc.Objects,
		TB:      tb,
	}
}
