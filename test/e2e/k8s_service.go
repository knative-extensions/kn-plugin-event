//go:build e2e
// +build e2e

package e2e

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/reconciler-test/pkg/feature"
)

// SendEventToKubeService returns a feature.Feature that verifies the kn-event
// can send to Kubernetes service.
func SendEventToKubeService() *feature.Feature {
	return SendEventFeature(kubeServiceSut{})
}

type kubeServiceSut struct{}

func (k kubeServiceSut) Name() string {
	return "KubeService"
}

func (k kubeServiceSut) Deploy(_ *feature.Feature, sinkName string) Sink {
	return sinkFn(func() string {
		return fmt.Sprintf("Service:%s:%s",
			corev1.SchemeGroupVersion, sinkName)
	})
}
