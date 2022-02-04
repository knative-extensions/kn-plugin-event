package reference

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/kmeta"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// FromKnativeService creates a reference of Service object.
func FromKnativeService(ctx context.Context, ksvc *servingv1.Service) corev1.ObjectReference {
	cp := ksvc.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}
