//go:build e2e

package e2e

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	messagingv1clientset "knative.dev/eventing/pkg/client/clientset/versioned/typed/messaging/v1"
	"knative.dev/kn-plugin-event/pkg/tests/reference"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/k8s"
)

// SendEventToChannel returns a feature.Feature that verifies the kn-event
// can send to Knative channel.
func SendEventToChannel() *feature.Feature {
	return SendEventFeature(channelSut{})
}

type channelSut struct{}

func (cs channelSut) Name() string {
	return "Channel"
}

func (cs channelSut) Deploy(f *feature.Feature, sinkName string) Sink {
	ch := channelSutImpl{sinkName}

	f.Setup("Deploy Channel", ch.step)

	return ch.sink()
}

type channelSutImpl struct {
	sinkName string
}

func (c channelSutImpl) step(ctx context.Context, t feature.T) {
	c.deployChannel(ctx, t)
	c.deploySubscription(ctx, t)
}

func (c channelSutImpl) sink() Sink {
	return sinkFn(func() string {
		return "channel:" + c.name()
	})
}

func (c channelSutImpl) name() string {
	return c.sinkName + "-channel"
}

func (c channelSutImpl) deployChannel(ctx context.Context, t feature.T) {
	env := environment.FromContext(ctx)
	ns := env.Namespace()
	rest := injection.GetConfig(ctx)
	channels := messagingv1clientset.NewForConfigOrDie(rest).Channels(ns)
	channel := &messagingv1.Channel{
		ObjectMeta: metav1.ObjectMeta{Name: c.name(), Namespace: ns},
		Spec: messagingv1.ChannelSpec{
			ChannelableSpec: eventingduckv1.ChannelableSpec{
				Delivery: deliverySpec(),
			},
		},
	}
	log := logging.FromContext(ctx).
		With(json("meta", channel.ObjectMeta))
	log.Info("Deploying Channel")
	if _, err := channels.Create(ctx, channel, metav1.CreateOptions{}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	ref := reference.FromChannel(ctx, channel)
	env.Reference(ref)
	k8s.WaitForReadyOrDoneOrFail(ctx, t, ref)
	log.Info("Channel is ready")
}

func (c channelSutImpl) deploySubscription(ctx context.Context, t feature.T) {
	env := environment.FromContext(ctx)
	ns := env.Namespace()
	rest := injection.GetConfig(ctx)
	subscriptions := messagingv1clientset.NewForConfigOrDie(rest).Subscriptions(ns)
	subscription := &messagingv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{Name: c.name(), Namespace: ns},
		Spec: messagingv1.SubscriptionSpec{
			Channel: *reference.ToKnative(reference.FromChannel(ctx, &messagingv1.Channel{
				ObjectMeta: metav1.ObjectMeta{Name: c.name()},
			})),
			Subscriber: &duckv1.Destination{
				Ref: reference.ToKnative(reference.FromKubeService(ctx, &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: c.sinkName},
				})),
			},
		},
	}
	log := logging.FromContext(ctx).
		With(json("meta", subscription.ObjectMeta))
	log.With(json("spec", subscription.Spec)).
		Info("Deploying Subscription")
	if _, err := subscriptions.Create(ctx, subscription, metav1.CreateOptions{}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	ref := reference.FromSubscription(ctx, subscription)
	env.Reference(ref)
	k8s.WaitForReadyOrDoneOrFail(ctx, t, ref)
	log.Info("Subscription is ready")
}
