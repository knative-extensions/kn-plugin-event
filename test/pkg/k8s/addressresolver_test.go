//go:build e2e

package k8s_test

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
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
	"knative.dev/pkg/logging"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestResolveAddress(t *testing.T) {
	plugintestpkg.WithClients(t, func(c plugintestpkg.ClientsContext) {
		for _, tc := range k8stest.ResolveAddressTestCases(c.KnTest.Namespace()) {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				log := zaptest.NewLogger(t, zaptest.Level(zapcore.InfoLevel))
				ctx := logging.WithLogger(context.TODO(), log.Sugar())
				k8stest.EnsureResolveAddress(ctx, t, tc, func() (k8s.Clients, func(tb testing.TB)) {
					deploy(ctx, t, tc, c.Clients)
					cleanup := func(tb testing.TB) { //nolint:thelper
						if tb.Failed() {
							tb.Logf("Skipping undeploy, because test '%s' failed", tb.Name())
							return
						}
						undeploy(ctx, tb, tc, c.Clients)
					}
					return c.Clients, cleanup
				})
			})
		}
	})
}

func deploy( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	tc k8stest.ResolveAddressTestCase,
	clients k8s.Clients,
) {
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

func undeploy( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	tc k8stest.ResolveAddressTestCase,
	clients k8s.Clients,
) {
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

func deployK8sService( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	service corev1.Service,
) {
	service.Status = corev1.ServiceStatus{}
	_, err := clients.Typed().CoreV1().Services(service.Namespace).
		Create(ctx, &service, metav1.CreateOptions{})
	assert.NilError(tb, err)
}

func undeployK8sService( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	service corev1.Service,
) {
	err := clients.Typed().CoreV1().Services(service.Namespace).
		Delete(ctx, service.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployKnService( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	service servingv1.Service,
) {
	service.Status = servingv1.ServiceStatus{}
	knclient := clientservingv1.NewKnServingClient(clients.Serving(), service.Namespace)
	err := knclient.CreateService(ctx, &service)
	assert.NilError(tb, err)
	err, _ = knclient.WaitForService(ctx, service.Name, clientservingv1.WaitConfig{
		Timeout:     time.Duration(2) * time.Minute,
		ErrorWindow: time.Duration(2) * time.Second,
	},
		clientwait.NoopMessageCallback())
	assert.NilError(tb, err)
}

func undeployKnService( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	service servingv1.Service,
) {
	err := clientservingv1.
		NewKnServingClient(clients.Serving(), service.Namespace).
		DeleteService(ctx, service.GetName(), time.Minute)
	assert.NilError(tb, err)
}

func deployBroker( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	broker eventingv1.Broker,
) {
	broker.Status = eventingv1.BrokerStatus{}
	knclient := clienteventingv1.NewKnEventingClient(clients.Eventing(),
		broker.Namespace)
	assert.NilError(tb, knclient.CreateBroker(ctx, &broker))
	waitForReady(ctx, tb, clients, &broker, time.Minute)
}

func undeployBroker( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	broker eventingv1.Broker,
) {
	err := clients.Eventing().Brokers(broker.Namespace).
		Delete(ctx, broker.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func deployChannel( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	channel messagingv1.Channel,
) {
	channel.Status = messagingv1.ChannelStatus{}
	knclient := clientmessagingv1.NewKnMessagingClient(clients.Messaging(),
		channel.Namespace).ChannelsClient()
	assert.NilError(tb, knclient.CreateChannel(ctx, &channel))
	waitForReady(ctx, tb, clients, &channel, time.Minute)
}

func undeployChannel( //nolint:thelper
	ctx context.Context,
	tb testing.TB,
	clients k8s.Clients,
	channel messagingv1.Channel,
) {
	err := clients.Messaging().Channels(channel.Namespace).
		Delete(ctx, channel.Name, metav1.DeleteOptions{})
	assert.NilError(tb, err)
}

func gvr(accessor kmeta.Accessor) schema.GroupVersionResource {
	gvk := accessor.GroupVersionKind()
	return apis.KindToResource(gvk)
}

func waitForReady(
	ctx context.Context,
	t poll.TestingT,
	clients k8s.Clients,
	accessor kmeta.Accessor,
	timeout time.Duration,
) {
	dynclient := clients.Dynamic()
	poll.WaitOn(t, isReady(ctx, dynclient, accessor),
		poll.WithTimeout(timeout), poll.WithDelay(time.Second))
}

func isReady(ctx context.Context, dynclient dynamic.Interface, accessor kmeta.Accessor) poll.Check {
	resources := dynclient.Resource(gvr(accessor)).
		Namespace(accessor.GetNamespace())
	return func(poll.LogT) poll.Result {
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
