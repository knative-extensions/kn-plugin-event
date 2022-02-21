//go:build e2e
// +build e2e

package k8s_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	clienteventingv1 "knative.dev/client/pkg/eventing/v1"
	clientmessagingv1 "knative.dev/client/pkg/messaging/v1"
	clientservingv1 "knative.dev/client/pkg/serving/v1"
	clientwait "knative.dev/client/pkg/wait"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	"knative.dev/kn-plugin-event/pkg/k8s"
	k8stest "knative.dev/kn-plugin-event/pkg/k8s/test"
	plugintestpkg "knative.dev/kn-plugin-event/test/pkg"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestResolveAddress(t *testing.T) {
	plugintestpkg.WithClients(t, func(c plugintestpkg.ClientsContext) {
		k8stest.ResolveAddressTestCases(c.KnTest.Namespace(), func(tc k8stest.ResolveAddressTestCase) {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				k8stest.EnsureResolveAddress(t, tc, func(ctx context.Context) (k8s.Clients, func(tb testing.TB)) {
					deploy(ctx, t, tc, c.Clients)
					cleanup := func(tb testing.TB) {
						tb.Helper()
						undeploy(ctx, tb, tc, c.Clients)
					}
					return c.Clients, cleanup
				})
			})
		})
	})
}

func deploy(
	ctx context.Context,
	tb testing.TB,
	tc k8stest.ResolveAddressTestCase,
	clients k8s.Clients,
) {
	tb.Helper()
	for _, object := range tc.Objects {
		switch v := object.(type) {
		case *servingv1.Service:
			deployKnService(ctx, tb, clients, *(v))
		case *corev1.Service:
			deployK8sService(ctx, tb, clients, *(v))
		case *eventingv1.Broker:
			deployBroker(ctx, tb, clients, *(v))
		case *messagingv1.Channel:
			deployChannel(ctx, tb, clients, *(v))
		default:
			tb.Fatalf("unsupported type: %#v", v)
		}
	}
}

func undeploy(
	ctx context.Context,
	tb testing.TB,
	tc k8stest.ResolveAddressTestCase,
	clients k8s.Clients,
) {
	tb.Helper()
	for _, object := range tc.Objects {
		switch v := object.(type) {
		case *servingv1.Service:
			undeployKnService(ctx, tb, clients, *(v))
		case *corev1.Service:
			undeployK8sService(ctx, tb, clients, *(v))
		case *eventingv1.Broker:
			undeployBroker(ctx, tb, clients, *(v))
		case *messagingv1.Channel:
			undeployChannel(ctx, tb, clients, *(v))
		default:
			tb.Fatalf("unsupported type: %#v", v)
		}
	}
}

func deployK8sService(ctx context.Context, tb testing.TB, clients k8s.Clients, service corev1.Service) {
	tb.Helper()
	service.Status = corev1.ServiceStatus{}
	_, err := clients.Typed().CoreV1().Services(service.Namespace).
		Create(ctx, &service, metav1.CreateOptions{})
	assert.NilError(tb, err)
}

func undeployK8sService(ctx context.Context, tb testing.TB, clients k8s.Clients, service corev1.Service) {
	tb.Helper()
	err := clients.Typed().CoreV1().Services(service.Namespace).
		Delete(ctx, service.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployKnService(ctx context.Context, tb testing.TB, clients k8s.Clients, service servingv1.Service) {
	tb.Helper()
	service.Status = servingv1.ServiceStatus{}
	knclient := clientservingv1.NewKnServingClient(clients.Serving(), service.Namespace)
	err := knclient.CreateService(ctx, &service)
	assert.NilError(tb, err)
	err, _ = knclient.WaitForService(ctx, service.Name, 2*time.Minute,
		clientwait.NoopMessageCallback())
	assert.NilError(tb, err)
}

func undeployKnService(ctx context.Context, tb testing.TB, clients k8s.Clients, service servingv1.Service) {
	tb.Helper()
	err := clientservingv1.
		NewKnServingClient(clients.Serving(), service.Namespace).
		DeleteService(ctx, service.GetName(), time.Minute)
	assert.NilError(tb, err)
}

func deployBroker(ctx context.Context, tb testing.TB, clients k8s.Clients, broker eventingv1.Broker) {
	tb.Helper()
	broker.Status = eventingv1.BrokerStatus{}
	knclient := clienteventingv1.NewKnEventingClient(clients.Eventing(),
		broker.Namespace)
	assert.NilError(tb, knclient.CreateBroker(ctx, &broker))
	assert.NilError(tb, waitForReady(ctx, clients, &broker, 30*time.Second))
}

func undeployBroker(ctx context.Context, tb testing.TB, clients k8s.Clients, broker eventingv1.Broker) {
	tb.Helper()
	err := clients.Eventing().Brokers(broker.Namespace).
		Delete(ctx, broker.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployChannel(ctx context.Context, tb testing.TB, clients k8s.Clients, channel messagingv1.Channel) {
	tb.Helper()
	channel.Status = messagingv1.ChannelStatus{}
	knclient := clientmessagingv1.NewKnMessagingClient(clients.Messaging(),
		channel.Namespace).ChannelsClient()
	assert.NilError(tb, knclient.CreateChannel(ctx, &channel))
	assert.NilError(tb, waitForReady(ctx, clients, &channel, 30*time.Second))
}

func undeployChannel(ctx context.Context, tb testing.TB, clients k8s.Clients, channel messagingv1.Channel) {
	tb.Helper()
	err := clients.Messaging().Channels(channel.Namespace).
		Delete(ctx, channel.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

type accessorWatchMaker struct {
	clients   k8s.Clients
	acccessor kmeta.Accessor
}

func (a accessorWatchMaker) watchMaker(
	ctx context.Context, name, _ string, _ time.Duration,
) (watch.Interface, error) {
	dynclient := a.clients.Dynamic().Resource(a.gvr()).
		Namespace(a.acccessor.GetNamespace())
	return dynclient.Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
	})
}

func (a accessorWatchMaker) gvr() schema.GroupVersionResource {
	gvk := a.acccessor.GroupVersionKind()
	return apis.KindToResource(gvk)
}

func waitForReady(ctx context.Context, clients k8s.Clients, acccessor kmeta.Accessor, timeout time.Duration) error {
	awm := accessorWatchMaker{clients: clients, acccessor: acccessor}
	ready := clientwait.NewWaitForReady(awm.gvr().Resource, awm.watchMaker, dynamicConditionExtractor)
	err, _ := ready.Wait(
		ctx,
		acccessor.GetName(),
		"0",
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
