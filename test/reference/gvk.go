package reference

import (
	"context"
	"errors"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	"knative.dev/pkg/logging"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

var (
	errNoKinds = errors.New("want len(kinds) > 0")
	once       sync.Once       //nolint:gochecknoglobals
	scheme     *runtime.Scheme //nolint:gochecknoglobals
)

func groupVersionKind(ctx context.Context, obj runtime.Object) schema.GroupVersionKind {
	log := logging.FromContext(ctx)
	once.Do(func() {
		scheme = runtime.NewScheme()
		builders := []runtime.SchemeBuilder{
			servingv1.SchemeBuilder,
			eventingv1.SchemeBuilder,
			messagingv1.SchemeBuilder,
			corev1.SchemeBuilder,
		}
		for _, builder := range builders {
			if err := builder.AddToScheme(scheme); err != nil {
				log.Fatal(err)
			}
		}
	})
	kinds, _, err := scheme.ObjectKinds(obj)
	if err != nil {
		log.Fatal(err)
	}
	if !(len(kinds) > 0) {
		log.Fatal(errNoKinds)
	}
	return kinds[0]
}
