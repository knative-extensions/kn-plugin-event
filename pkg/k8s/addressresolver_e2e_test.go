// +build e2e

package k8s_test

import (
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clienttest "knative.dev/client/lib/test"
	clientservingv1 "knative.dev/client/pkg/serving/v1"
	clientwait "knative.dev/client/pkg/wait"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/pkg/apis"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestResolveAddressE2E(t *testing.T) {
	clients, err := k8s.CreateKubeClient(&event.Properties{})
	if err != nil && errors.Is(err, k8s.ErrNoKubernetesConnection) {
		t.Skip("AUTO-SKIP:", err)
	} else {
		assert.NilError(t, err)
	}
	t.Parallel()
	it, err := clienttest.NewKnTest()
	assert.NilError(t, err)
	t.Cleanup(func() {
		assert.NilError(t, it.Teardown())
	})

	resolveAddressTestCases(it.Namespace(), func(tc resolveAddressTestCase) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			performResolveAddressTest(t, tc, func(tr resolveAddressTestResources) (k8s.Clients, func(tb testing.TB)) {
				deployK8sServices(t, tr, clients)
				deployKnServices(t, tr, clients)
				return clients, func(tb testing.TB) {
					tb.Helper()
					undeployK8sServices(tb, tr, clients)
					undeployKnServices(tb, tr, clients)
				}
			})
		})
	})
}

func deployK8sServices(tb testing.TB, tr resolveAddressTestResources, clients k8s.Clients) {
	tb.Helper()
	for _, service := range tr.k8sServices {
		//goland:noinspection GoShadowedVar
		service := service
		service.Status = corev1.ServiceStatus{}
		_, err := clients.Typed().CoreV1().Services(service.Namespace).
			Create(clients.Context(), &service, metav1.CreateOptions{})
		assert.NilError(tb, err)
	}
}

func undeployK8sServices(tb testing.TB, tr resolveAddressTestResources, clients k8s.Clients) {
	tb.Helper()
	for _, service := range tr.k8sServices {
		err := clients.Typed().CoreV1().Services(service.Namespace).
			Delete(clients.Context(), service.Name, metav1.DeleteOptions{})
		assert.NilError(tb, err)
	}
}

func servingGVK() schema.GroupVersionResource {
	return apis.KindToResource(
		servingv1.SchemeGroupVersion.WithKind("Service"),
	)
}

func deployKnServices(tb testing.TB, tr resolveAddressTestResources, clients k8s.Clients) {
	tb.Helper()
	for _, service := range tr.knServices {
		//goland:noinspection GoShadowedVar
		service := service
		service.Status = servingv1.ServiceStatus{}
		ctx := clients.Context()
		knclient := clientservingv1.NewKnServingClient(clients.Serving(), service.Namespace)
		err := knclient.CreateService(ctx, &service)
		assert.NilError(tb, err)
		err, _ = knclient.WaitForService(ctx, service.Name, 2*time.Minute,
			clientwait.NoopMessageCallback())
		assert.NilError(tb, err)
	}
}

func undeployKnServices(tb testing.TB, tr resolveAddressTestResources, clients k8s.Clients) {
	tb.Helper()
	for _, service := range tr.knServices {
		err := clients.Dynamic().Resource(servingGVK()).
			Namespace(service.Namespace).
			Delete(clients.Context(), service.Name, metav1.DeleteOptions{})
		assert.NilError(tb, err)
	}
}
