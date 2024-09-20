//go:build e2e

package e2e_test

import (
	"os"
	"testing"

	// ref.: https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"knative.dev/reconciler-test/pkg/environment"
)

// Global is the singleton instance of GlobalEnvironment. It is used to parse
// the testing config for the test run. The config will specify the cluster
// config as well as the parsing level and state flags.
var global environment.GlobalEnvironment //nolint:gochecknoglobals

// TestMain is the first entry point for `go test`.
func TestMain(m *testing.M) {
	global = environment.NewStandardGlobalEnvironment()

	// Run the tests.
	os.Exit(m.Run())
}
