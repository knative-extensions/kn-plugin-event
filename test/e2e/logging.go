//go:build e2e
// +build e2e

package e2e

import (
	"go.uber.org/zap"
	"knative.dev/kn-plugin-event/pkg/tests/logging"
)

func json(name string, o interface{}) zap.Field {
	return logging.JSON(name, o)
}
