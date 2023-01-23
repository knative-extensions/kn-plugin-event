package pkg

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
	clienttest "knative.dev/client-pkg/pkg/util/test"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	plugintest "knative.dev/kn-plugin-event/test"
)

// TestContext holds a test context.
type TestContext struct {
	testing.TB
	*clienttest.KnTest
}

// ClientsContext holds a Kubernetes clients context.
type ClientsContext struct {
	k8s.Clients
	*TestContext
}

// WithClients runs provided handler within a ClientsContext, so that the
// handler has access to configured k8s.Clients.
func WithClients(tb testing.TB, handler func(c ClientsContext)) {
	tb.Helper()
	plugintest.MaybeSkip(tb)
	clients, err := k8s.CreateKubeClient(&event.Properties{})
	if err != nil && errors.Is(err, k8s.ErrNoKubernetesConnection) {
		tb.Skip("AUTO-SKIP:", err)
	} else {
		assert.NilError(tb, err)
	}
	WithKnTest(tb, func(c *TestContext) {
		handler(ClientsContext{
			TestContext: c, Clients: clients,
		})
	})
}

// WithKnTest runs handler within a TestContext, so that the handler has access
// to clienttest.KnTest.
func WithKnTest(tb testing.TB, handler func(c *TestContext)) {
	tb.Helper()
	plugintest.MaybeSkip(tb)
	it, err := clienttest.NewKnTest()
	assert.NilError(tb, err)
	tb.Cleanup(func() {
		if tb.Failed() {
			tb.Logf("Skipping '%s' namespace teardown because '%s' test is failing",
				it.Namespace(), tb.Name())
			return
		}
		assert.NilError(tb, it.Teardown())
	})
	handler(&TestContext{
		TB: tb, KnTest: it,
	})
}
