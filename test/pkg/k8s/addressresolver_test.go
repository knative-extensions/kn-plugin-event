//go:build e2e
// +build e2e

package k8s_test

import (
	"context"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/poll"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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
				k8stest.EnsureResolveAddress(t, tc, func() (k8s.Clients, func(tb testing.TB)) {
					deploy(t, tc, c.Clients)
					cleanup := func(tb testing.TB) { // nolint:thelper
						if tb.Failed() {
							tb.Logf("Skipping undeploy, because test '%s' failed", tb.Name())
							return
						}
						undeploy(tb, tc, c.Clients)
					}
					return c.Clients, cleanup
				})
			})
		})
	})
}

func deploy(tb testing.TB, tc k8stest.ResolveAddressTestCase, clients k8s.Clients) { // nolint:thelper
	for _, object := range tc.Objects {
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

func undeploy(tb testing.TB, tc k8stest.ResolveAddressTestCase, clients k8s.Clients) { // nolint:thelper
	for _, object := range tc.Objects {
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

func deployK8sService(tb testing.TB, clients k8s.Clients, service corev1.Service) { // nolint:thelper
	service.Status = corev1.ServiceStatus{}
	_, err := clients.Typed().CoreV1().Services(service.Namespace).
		Create(clients.Context(), &service, metav1.CreateOptions{})
	assert.NilError(tb, err)
}

func undeployK8sService(tb testing.TB, clients k8s.Clients, service corev1.Service) { // nolint:thelper
	err := clients.Typed().CoreV1().Services(service.Namespace).
		Delete(clients.Context(), service.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployKnService(tb testing.TB, clients k8s.Clients, service servingv1.Service) { // nolint:thelper
	service.Status = servingv1.ServiceStatus{}
	ctx := clients.Context()
	knclient := clientservingv1.NewKnServingClient(clients.Serving(), service.Namespace)
	err := knclient.CreateService(ctx, &service)
	assert.NilError(tb, err)
	err, _ = knclient.WaitForService(ctx, service.Name, 2*time.Minute,
		clientwait.NoopMessageCallback())
	assert.NilError(tb, err)
}

func undeployKnService(tb testing.TB, clients k8s.Clients, service servingv1.Service) { // nolint:thelper
	err := clientservingv1.
		NewKnServingClient(clients.Serving(), service.Namespace).
		DeleteService(clients.Context(), service.GetName(), time.Minute)
	assert.NilError(tb, err)
}

func deployBroker(tb testing.TB, clients k8s.Clients, broker eventingv1.Broker) { // nolint:thelper
	broker.Status = eventingv1.BrokerStatus{}
	ctx := clients.Context()
	knclient := clienteventingv1.NewKnEventingClient(clients.Eventing(),
		broker.Namespace)
	assert.NilError(tb, knclient.CreateBroker(ctx, &broker))
	waitForReady(tb, clients, &broker, time.Minute)
}

func undeployBroker(tb testing.TB, clients k8s.Clients, broker eventingv1.Broker) { // nolint:thelper
	err := clients.Eventing().Brokers(broker.Namespace).
		Delete(clients.Context(), broker.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployChannel(tb testing.TB, clients k8s.Clients, channel messagingv1.Channel) { // nolint:thelper
	channel.Status = messagingv1.ChannelStatus{}
	knclient := clientmessagingv1.NewKnMessagingClient(clients.Messaging(),
		channel.Namespace).ChannelsClient()
	assert.NilError(tb, knclient.CreateChannel(clients.Context(), &channel))
	waitForReady(tb, clients, &channel, time.Minute)
}

func undeployChannel(tb testing.TB, clients k8s.Clients, channel messagingv1.Channel) { // nolint:thelper
	err := clients.Messaging().Channels(channel.Namespace).
		Delete(clients.Context(), channel.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func gvr(accessor kmeta.Accessor) schema.GroupVersionResource {
	gvk := accessor.GroupVersionKind()
	return apis.KindToResource(gvk)
}

func waitForReady(t poll.TestingT, clients k8s.Clients, accessor kmeta.Accessor, timeout time.Duration) {
	ctx := clients.Context()
	dynclient := clients.Dynamic()
	poll.WaitOn(t, isReady(ctx, dynclient, accessor),
		poll.WithTimeout(timeout), poll.WithDelay(time.Second))
}

func isReady(ctx context.Context, dynclient dynamic.Interface, accessor kmeta.Accessor) poll.Check {
	resources := dynclient.Resource(gvr(accessor)).
		Namespace(accessor.GetNamespace())
	return func(t poll.LogT) poll.Result {
		res, err := resources.Get(ctx, accessor.GetName(), metav1.GetOptions{})
		if err != nil {
			return poll.Error(err)
		}
		kres, err := toKResource(res)
		if err != nil {
			return poll.Error(err)
		}
		for _, cond := range kres.Status.Conditions {
			if cond.Type == apis.ConditionReady {
				if cond.Status == corev1.ConditionTrue {
					return poll.Success()
				}
				return poll.Continue(
					"%s named '%s' in namespace '%s' is not ready '%s', reason '%s'",
					accessor.GroupVersionKind(), accessor.GetName(),
					accessor.GetNamespace(), cond.Status, cond.Reason)
			}
		}
		return poll.Continue(
			"%s named '%s' in namespace '%s' does not have ready condition",
			accessor.GroupVersionKind(), accessor.GetName(),
			accessor.GetNamespace())
	}
}

func toKResource(obj runtime.Object) (*duckv1.KResource, error) {
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
	return &kresource, nil
}
