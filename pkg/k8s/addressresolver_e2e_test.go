// +build e2e

package k8s_test

import (
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
	"gotest.tools/v3/assert"
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
	t.Parallel()
	it, err := clienttest.NewKnTest()
	assert.NilError(t, err)
	t.Cleanup(func() {
		assert.NilError(t, it.Teardown())
	})

	knconfig, err := homedir.Expand("~/.config/kn/config.yaml")
	assert.NilError(t, err)
	kubeconfig, err := homedir.Expand("~/.kube/config")
	assert.NilError(t, err)
	props := &event.Properties{
		KnPluginOptions: event.KnPluginOptions{
			KnConfig:   knconfig,
			Kubeconfig: kubeconfig,
		},
	}
	clients, err := k8s.CreateKubeClient(props)
	assert.NilError(t, err)

	testCases := []resolveAddressTestCase{
		k8sService(it.Namespace()),
		knServiceClusterLocal(it.Namespace()),
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			performResolveAddressTest(t, tc, func(tr resolveAddressTestResources) (k8s.Clients, func(t *testing.T)) {
				deployK8sServices(t, tr, clients)
				deployKnServices(t, tr, clients)
				return clients, func(t *testing.T) {
					undeployK8sServices(t, tr, clients)
					undeployKnServices(t, tr, clients)
				}
			})
		})
	}
}

func deployK8sServices(t *testing.T, tr resolveAddressTestResources, clients k8s.Clients) {
	for i := range tr.k8sServices {
		service := tr.k8sServices[i]
		_, err := clients.Typed().CoreV1().Services(service.Namespace).
			Create(clients.Context(), &service, metav1.CreateOptions{})
		assert.NilError(t, err)
	}
}

func undeployK8sServices(t *testing.T, tr resolveAddressTestResources, clients k8s.Clients) {
	for _, service := range tr.k8sServices {
		err := clients.Typed().CoreV1().Services(service.Namespace).
			Delete(clients.Context(), service.Name, metav1.DeleteOptions{})
		assert.NilError(t, err)
	}
}

func servingGVK() schema.GroupVersionResource {
	return apis.KindToResource(
		servingv1.SchemeGroupVersion.WithKind("Service"),
	)
}

func deployKnServices(t *testing.T, tr resolveAddressTestResources, clients k8s.Clients) {
	for i := range tr.knServices {
		service := tr.knServices[i]
		ctx := clients.Context()
		knclient := clientservingv1.NewKnServingClient(clients.Serving(), service.Namespace)
		err := knclient.CreateService(ctx, &service)
		assert.NilError(t, err)
		err, _ = knclient.WaitForService(ctx, service.Name, 2 * time.Minute,
			clientwait.NoopMessageCallback())
		assert.NilError(t, err)
	}
}

func undeployKnServices(t *testing.T, tr resolveAddressTestResources, clients k8s.Clients) {
	for _, service := range tr.knServices {
		err := clients.Dynamic().Resource(servingGVK()).
			Namespace(service.Namespace).
			Delete(clients.Context(), service.Name, metav1.DeleteOptions{})
		assert.NilError(t, err)
	}
}
