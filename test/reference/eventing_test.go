package reference_test

import (
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	"knative.dev/kn-plugin-event/test/reference"
)

func TestFromBroker(t *testing.T) {
	ctx := testContext(t)
	broker := &eventingv1.Broker{
		ObjectMeta: meta("broker", "foo"),
	}

	got := reference.FromBroker(ctx, broker)
	want := corev1.ObjectReference{
		Kind:       "Broker",
		Namespace:  broker.Namespace,
		Name:       broker.Name,
		APIVersion: eventingv1.SchemeGroupVersion.String(),
	}
	assert.Equal(t, want, got)
}

func TestTrigger(t *testing.T) {
	ctx := testContext(t)
	trigger := &eventingv1.Trigger{
		ObjectMeta: meta("trigger", "foo"),
	}

	got := reference.FromTrigger(ctx, trigger)
	want := corev1.ObjectReference{
		Kind:       "Trigger",
		Namespace:  trigger.Namespace,
		Name:       trigger.Name,
		APIVersion: eventingv1.SchemeGroupVersion.String(),
	}
	assert.Equal(t, want, got)
}

func TestChannel(t *testing.T) {
	ctx := testContext(t)
	channel := &messagingv1.Channel{
		ObjectMeta: meta("channel", "fizz"),
	}

	got := reference.FromChannel(ctx, channel)
	want := corev1.ObjectReference{
		Kind:       "Channel",
		Namespace:  channel.Namespace,
		Name:       channel.Name,
		APIVersion: messagingv1.SchemeGroupVersion.String(),
	}
	assert.Equal(t, want, got)
}

func TestSubscription(t *testing.T) {
	ctx := testContext(t)
	subscription := &messagingv1.Subscription{
		ObjectMeta: meta("subscription", "bazz"),
	}

	got := reference.FromSubscription(ctx, subscription)
	want := corev1.ObjectReference{
		Kind:       "Subscription",
		Namespace:  subscription.Namespace,
		Name:       subscription.Name,
		APIVersion: messagingv1.SchemeGroupVersion.String(),
	}
	assert.Equal(t, want, got)
}
