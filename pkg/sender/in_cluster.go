package sender

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"knative.dev/kn-plugin-event/pkg/cli/ics"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

const idLength = 16

type inClusterSender struct {
	addressable     *event.AddressableSpec
	jobRunner       k8s.JobRunner
	addressResolver k8s.ReferenceAddressResolver
}

func (i *inClusterSender) Send(ce cloudevents.Event) error {
	url, err := i.addressResolver.ResolveAddress(
		i.addressable.Reference, i.addressable.URI,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", k8s.ErrInvalidReference, err)
	}
	kevent, err := ics.Encode(ce)
	if err != nil {
		return fmt.Errorf("%w: %w", ics.ErrCouldntEncode, err)
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      newJobName(),
			Namespace: i.addressable.SenderNamespace,
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
	err = i.jobRunner.Run(job)
	if err != nil {
		return fmt.Errorf("%w: %w", ics.ErrCantSendWithICS, err)
	}
	return nil
}

func newJobName() string {
	id := rand.String(idLength)
	return fmt.Sprintf("kn-event-sender-%s", id)
}
