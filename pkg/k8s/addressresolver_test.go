package k8s_test

import (
	"testing"

	clienttest "knative.dev/client/pkg/util/test"
	"knative.dev/kn-plugin-event/pkg/k8s"
	k8stest "knative.dev/kn-plugin-event/pkg/k8s/test"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestResolveAddress(t *testing.T) {
	ns := clienttest.NextNamespace()
	k8stest.ResolveAddressTestCases(ns, func(tc k8stest.ResolveAddressTestCase) {
		t.Run(tc.Name, func(t *testing.T) {
			k8stest.EnsureResolveAddress(t, tc, func() (k8s.Clients, func(tb testing.TB)) {
				return fakeClients(t, tc), noCleanup
			})
		})
	})
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
