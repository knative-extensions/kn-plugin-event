package test

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
	clienttest "knative.dev/client/lib/test"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
)

type Context struct {
	testing.TB
	*clienttest.KnTest
}

type ClientsContext struct {
	k8s.Clients
	*Context
}

func WithClients(tb testing.TB, handler func(c ClientsContext)) {
	tb.Helper()
	maybeSkip(tb, "clients")
	clients, err := k8s.CreateKubeClient(&event.Properties{})
	if err != nil && errors.Is(err, k8s.ErrNoKubernetesConnection) {
		tb.Skip("AUTO-SKIP:", err)
	} else {
		assert.NilError(tb, err)
	}
	WithKnTest(tb, func(c *Context) {
		handler(ClientsContext{
			Context: c, Clients: clients,
		})
	})
}

func WithKnTest(tb testing.TB, handler func(c *Context)) {
	tb.Helper()
	maybeSkip(tb, "kn")
	it, err := clienttest.NewKnTest()
	assert.NilError(tb, err)
	tb.Cleanup(func() {
		assert.NilError(tb, it.Teardown())
	})
	handler(&Context{
		TB: tb, KnTest: it,
	})
}

func maybeSkip(tb testing.TB, thing string) {
	if testing.Short() {
		tb.Skipf("Short flag is set. Skipping %s-test.", thing)
	}
}
