// +build e2e

package k8s_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	clienttest "knative.dev/client/lib/test"
	clienteventingv1 "knative.dev/client/pkg/eventing/v1"
	clientmessagingv1 "knative.dev/client/pkg/messaging/v1"
	clientservingv1 "knative.dev/client/pkg/serving/v1"
	clientwait "knative.dev/client/pkg/wait"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
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
			performResolveAddressTest(t, tc, func() (k8s.Clients, func(tb testing.TB)) {
				deploy(t, tc, clients)
				cleanup := func(tb testing.TB) {
					tb.Helper()
					undeploy(tb, tc, clients)
				}
				return clients, cleanup
			})
		})
	})
}

func deploy(tb testing.TB, tc resolveAddressTestCase, clients k8s.Clients) {
	tb.Helper()
	for _, object := range tc.objects {
		switch v := object.(type) {
		case *servingv1.Service:
			deployKnService(tb, clients, *(v))
		case *corev1.Service:
			deployK8sService(tb, clients, *(v))
		case *eventingv1.Broker:
			deployBroker(tb, clients, *(v))
		case *messagingv1.Channel:
			deployChannel(tb, clients, *(v))
		default:
			tb.Fatalf("unsupported type: %#v", v)
		}
	}
}

func undeploy(tb testing.TB, tc resolveAddressTestCase, clients k8s.Clients) {
	tb.Helper()
	for _, object := range tc.objects {
		switch v := object.(type) {
		case *servingv1.Service:
			undeployKnService(tb, clients, *(v))
		case *corev1.Service:
			undeployK8sService(tb, clients, *(v))
		case *eventingv1.Broker:
			undeployBroker(tb, clients, *(v))
		case *messagingv1.Channel:
			undeployChannel(tb, clients, *(v))
		default:
			tb.Fatalf("unsupported type: %#v", v)
		}
	}
}

func deployK8sService(tb testing.TB, clients k8s.Clients, service corev1.Service) {
	tb.Helper()
	service.Status = corev1.ServiceStatus{}
	_, err := clients.Typed().CoreV1().Services(service.Namespace).
		Create(clients.Context(), &service, metav1.CreateOptions{})
	assert.NilError(tb, err)
}

func undeployK8sService(tb testing.TB, clients k8s.Clients, service corev1.Service) {
	tb.Helper()
	err := clients.Typed().CoreV1().Services(service.Namespace).
		Delete(clients.Context(), service.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployKnService(tb testing.TB, clients k8s.Clients, service servingv1.Service) {
	tb.Helper()
	service.Status = servingv1.ServiceStatus{}
	ctx := clients.Context()
	knclient := clientservingv1.NewKnServingClient(clients.Serving(), service.Namespace)
	err := knclient.CreateService(ctx, &service)
	assert.NilError(tb, err)
	err, _ = knclient.WaitForService(ctx, service.Name, 2*time.Minute,
		clientwait.NoopMessageCallback())
	assert.NilError(tb, err)
}

func undeployKnService(tb testing.TB, clients k8s.Clients, service servingv1.Service) {
	tb.Helper()
	err := clientservingv1.
		NewKnServingClient(clients.Serving(), service.Namespace).
		DeleteService(clients.Context(), service.GetName(), time.Minute)
	assert.NilError(tb, err)
}

func deployBroker(tb testing.TB, clients k8s.Clients, broker eventingv1.Broker) {
	tb.Helper()
	broker.Status = eventingv1.BrokerStatus{}
	ctx := clients.Context()
	knclient := clienteventingv1.NewKnEventingClient(clients.Eventing(),
		broker.Namespace)
	assert.NilError(tb, knclient.CreateBroker(ctx, &broker))
	assert.NilError(tb, waitForReady(clients, &broker, 30*time.Second))
}

func undeployBroker(tb testing.TB, clients k8s.Clients, broker eventingv1.Broker) {
	tb.Helper()
	err := clients.Eventing().Brokers(broker.Namespace).
		Delete(clients.Context(), broker.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployChannel(tb testing.TB, clients k8s.Clients, channel messagingv1.Channel) {
	tb.Helper()
	channel.Status = messagingv1.ChannelStatus{}
	knclient := clientmessagingv1.NewKnMessagingClient(clients.Messaging(),
		channel.Namespace).ChannelsClient()
	assert.NilError(tb, knclient.CreateChannel(clients.Context(), &channel))
	assert.NilError(tb, waitForReady(clients, &channel, 30*time.Second))
}

func undeployChannel(tb testing.TB, clients k8s.Clients, channel messagingv1.Channel) {
	tb.Helper()
	err := clients.Messaging().Channels(channel.Namespace).
		Delete(clients.Context(), channel.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func waitForReady(clients k8s.Clients, acccessor kmeta.Accessor, timeout time.Duration) error {
	gvk := acccessor.GroupVersionKind()
	gvr := apis.KindToResource(gvk)
	dynclient := clients.Dynamic().Resource(gvr).Namespace(acccessor.GetNamespace())
	ctx, cancel := context.WithTimeout(clients.Context(), timeout)
	defer cancel()
	watcher, err := dynclient.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	ready := clientwait.NewWaitForReady(gvr.Resource, dynamicConditionExtractor)
	err, _ = ready.Wait(
		ctx,
		watcher,
		acccessor.GetName(),
		clientwait.Options{Timeout: &timeout},
		clientwait.NoopMessageCallback(),
	)
	return err
}

func dynamicConditionExtractor(obj runtime.Object) (apis.Conditions, error) {
	un, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, k8s.ErrUnexcpected
	}
	kresource := duckv1.KResource{}
	err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(un.UnstructuredContent(), &kresource)
	if err != nil {
		return nil, err
	}
	return kresource.GetStatus().GetConditions(), nil
}
