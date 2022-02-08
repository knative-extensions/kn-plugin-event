package reference

import (
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ToKnative converts corev1.ObjectReference to duckv1.KReference.
func ToKnative(ref corev1.ObjectReference) *duckv1.KReference {
	return &duckv1.KReference{
		Kind:       ref.Kind,
		Namespace:  ref.Namespace,
		Name:       ref.Name,
		APIVersion: ref.APIVersion,
	}
}
