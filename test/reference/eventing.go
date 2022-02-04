package reference

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	"knative.dev/pkg/kmeta"
)

// FromBroker creates a reference of a Broker object.
func FromBroker(ctx context.Context, broker *eventingv1.Broker) corev1.ObjectReference {
	cp := broker.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}

// FromTrigger creates a reference of a Trigger object.
func FromTrigger(ctx context.Context, trigger *eventingv1.Trigger) corev1.ObjectReference {
	cp := trigger.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}

// FromChannel creates a reference of a Channel object.
func FromChannel(ctx context.Context, channel *messagingv1.Channel) corev1.ObjectReference {
	cp := channel.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}

// FromSubscription creates a reference of a Subscription object.
func FromSubscription(ctx context.Context, subscription *messagingv1.Subscription) corev1.ObjectReference {
	cp := subscription.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}
