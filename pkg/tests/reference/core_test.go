package reference_test

import (
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/kn-plugin-event/pkg/tests/reference"
)

func TestConfigMap(t *testing.T) {
	ctx := testContext(t)
	cm := &corev1.ConfigMap{
		ObjectMeta: meta("config-map", "bar"),
	}

	got := reference.FromConfigMap(ctx, cm)
	want := corev1.ObjectReference{
		Kind:       "ConfigMap",
		Namespace:  cm.Namespace,
		Name:       cm.Name,
		APIVersion: "v1",
	}
	assert.Equal(t, want, got)
}
