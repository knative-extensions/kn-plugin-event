package reference_test

import (
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	"knative.dev/kn-plugin-event/test/reference"
)

func TestBroker(t *testing.T) {
	ctx := testContext(t)
	broker := &eventingv1.Broker{
		ObjectMeta: meta("broker", "foo"),
	}

	got := reference.Broker(ctx, broker)
	want := corev1.ObjectReference{
		Kind:       "Broker",
		Namespace:  broker.Namespace,
		Name:       broker.Name,
		APIVersion: "eventing.knative.dev/v1",
	}
	assert.Equal(t, want, got)
}

func TestTrigger(t *testing.T) {
	ctx := testContext(t)
	trigger := &eventingv1.Trigger{
		ObjectMeta: meta("trigger", "foo"),
	}

	got := reference.Trigger(ctx, trigger)
	want := corev1.ObjectReference{
		Kind:       "Trigger",
		Namespace:  trigger.Namespace,
		Name:       trigger.Name,
		APIVersion: "eventing.knative.dev/v1",
	}
	assert.Equal(t, want, got)
}
