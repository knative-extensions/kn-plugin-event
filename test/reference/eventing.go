package reference

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	"knative.dev/pkg/kmeta"
)

func Broker(ctx context.Context, broker *eventingv1.Broker) corev1.ObjectReference {
	cp := broker.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}

func Trigger(ctx context.Context, trigger *eventingv1.Trigger) corev1.ObjectReference {
	cp := trigger.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}
