package sender

import (
	"context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"knative.dev/kn-plugin-event/pkg/event"
	ics2 "knative.dev/kn-plugin-event/pkg/ics"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

const idLength = 16

type inClusterSender struct {
	namespace       string
	target          *event.Target
	jobRunner       k8s.JobRunner
	addressResolver k8s.ReferenceAddressResolver
}

func (i *inClusterSender) Send(ctx context.Context, ce cloudevents.Event) error {
	url, err := i.addressResolver.ResolveAddress(
		ctx, i.target.Reference, i.target.RelativeURI,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", k8s.ErrInvalidReference, err)
	}
	kevent, err := ics2.Encode(ce)
	if err != nil {
		return fmt.Errorf("%w: %w", ics2.ErrCouldntEncode, err)
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      newJobName(),
			Namespace: i.namespace,
			Labels: map[string]string{
				"event-id": ce.ID(),
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{{
						Name:  "kn-event-sender",
						Image: metadata.ResolveImage(),
						Env: []corev1.EnvVar{{
							Name:  "K_SINK",
							Value: url.String(),
						}, {
							Name:  "K_EVENT",
							Value: kevent,
						}},
					}},
				},
			},
		},
	}
	err = i.jobRunner.Run(ctx, job)
	if err != nil {
		return fmt.Errorf("%w: %w", ics2.ErrCantSendWithICS, err)
	}
	return nil
}

func newJobName() string {
	id := rand.String(idLength)
	return "kn-event-sender-" + id
}
