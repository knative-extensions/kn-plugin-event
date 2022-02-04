package reference

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/kmeta"
)

// FromConfigMap create a reference of a ConfigMap object.
func FromConfigMap(ctx context.Context, cm *corev1.ConfigMap) corev1.ObjectReference {
	cp := cm.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}

// FromKubeService creates a reference of a Kubernetes Service object.
func FromKubeService(ctx context.Context, service *corev1.Service) corev1.ObjectReference {
	cp := service.DeepCopy()
	kind := groupVersionKind(ctx, cp)
	cp.APIVersion, cp.Kind = kind.ToAPIVersionAndKind()
	return kmeta.ObjectReference(cp)
}
