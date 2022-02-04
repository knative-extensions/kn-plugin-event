//go:build e2e
// +build e2e

package e2e

import (
	"knative.dev/reconciler-test/pkg/feature"
)

// SendEventToKubeService returns a feature.Feature that verifies the kn-event
// can send to Kubernetes service.
func SendEventToKubeService() *feature.Feature {
	return sendEventFeature("ToKubeService", sendEventOptions{
		sink: sinkFormat("Service:v1:%s"),
	})
}
