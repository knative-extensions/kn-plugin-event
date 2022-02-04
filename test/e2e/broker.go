//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	eventingv1clientset "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1"
	"knative.dev/kn-plugin-event/test/reference"
	"knative.dev/pkg/injection"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/k8s"
	"knative.dev/reconciler-test/resources/svc"
)

// SendEventToBroker returns a feature.Feature that verifies the kn-event
// can send to Knative broker.
func SendEventToBroker() *feature.Feature {
	return SendEventFeature(brokerSut{})
}

type brokerSut struct{}

func (bs brokerSut) Name() string {
	return "Broker"
}

func (bs brokerSut) Deploy(f *feature.Feature, sinkName string) Sink {
	b := brokerSutImpl{sinkName}
	f.Setup("deploy Broker", b.step)
	return b.sink()
}

type brokerSutImpl struct {
	sinkName string
}

func (b brokerSutImpl) step(ctx context.Context, t feature.T) {
	b.deployBroker(ctx, t)
	b.deployTrigger(ctx, t)
}

func (b brokerSutImpl) sink() Sink {
	return sinkFn(func() string {
		return fmt.Sprintf("Broker:%s:%s",
			eventingv1.SchemeGroupVersion, b.name())
	})
}

func (b brokerSutImpl) name() string {
	return b.sinkName + "-broker"
}

func (b brokerSutImpl) deployBroker(ctx context.Context, t feature.T) {
	env := environment.FromContext(ctx)
	ns := env.Namespace()
	rest := injection.GetConfig(ctx)
	brokers := eventingv1clientset.NewForConfigOrDie(rest).Brokers(ns)
	broker := &eventingv1.Broker{
		ObjectMeta: metav1.ObjectMeta{Name: b.name(), Namespace: ns},
		Spec: eventingv1.BrokerSpec{
			Delivery: deliverySpec(),
		},
	}
	if _, err := brokers.Create(ctx, broker, metav1.CreateOptions{}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	ref := reference.FromBroker(ctx, broker)
	env.Reference(ref)
	k8s.WaitForReadyOrDoneOrFail(ctx, t, ref)
}

func (b brokerSutImpl) deployTrigger(ctx context.Context, t feature.T) {
	env := environment.FromContext(ctx)
	ns := env.Namespace()
	rest := injection.GetConfig(ctx)
	triggers := eventingv1clientset.NewForConfigOrDie(rest).Triggers(ns)
	sinkRef := svc.AsDestinationRef(b.sinkName)
	sinkRef.Ref.Namespace = ns
	trigger := &eventingv1.Trigger{
		ObjectMeta: metav1.ObjectMeta{Name: b.name(), Namespace: ns},
		Spec: eventingv1.TriggerSpec{
			Broker:     b.name(),
			Subscriber: *sinkRef,
		},
	}
	if _, err := triggers.Create(ctx, trigger, metav1.CreateOptions{}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	ref := reference.FromTrigger(ctx, trigger)
	env.Reference(ref)
	k8s.WaitForReadyOrDoneOrFail(ctx, t, ref)
}

func deliverySpec() *eventingduckv1.DeliverySpec {
	const retryCount = 12
	var (
		retryCount32  = int32(retryCount)
		backoffPolicy = eventingduckv1.BackoffPolicyExponential
		backoffDelay  = "PT1S"
	)
	return &eventingduckv1.DeliverySpec{
		Retry:         &retryCount32,
		BackoffPolicy: &backoffPolicy,
		BackoffDelay:  &backoffDelay,
	}
}
