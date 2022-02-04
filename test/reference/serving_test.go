package reference_test

import (
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/kn-plugin-event/test/reference"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestKnativeService(t *testing.T) {
	ctx := testContext(t)
	service := &servingv1.Service{
		ObjectMeta: meta("ksvc", "foo"),
	}

	got := reference.FromKnativeService(ctx, service)
	want := corev1.ObjectReference{
		Kind:       "Service",
		Namespace:  service.Namespace,
		Name:       service.Name,
		APIVersion: servingv1.SchemeGroupVersion.String(),
	}
	assert.Equal(t, want, got)
}
