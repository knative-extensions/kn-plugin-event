package reference

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/kmeta"
)

func ConfigMap(ctx context.Context, cm *corev1.ConfigMap) corev1.ObjectReference {
	cp := cm.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}
