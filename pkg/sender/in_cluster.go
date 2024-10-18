package sender

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"knative.dev/kn-plugin-event/pkg/errors"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/ics"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/metadata"
	"knative.dev/pkg/ptr"
)

const (
	idLength        = 6
	defaultRetries  = 3
	defaultDeadline = 5 // seconds
)

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
		return errors.Wrap(err, k8s.ErrInvalidReference)
	}
	kevent, err := ics.Encode(ce)
	if err != nil {
		return errors.Wrap(err, ics.ErrCouldntEncode)
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
			BackoffLimit: ptr.Int32(defaultRetries),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:         corev1.RestartPolicyNever,
					ActiveDeadlineSeconds: ptr.Int64(defaultDeadline),
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
		return errors.Wrap(err, ics.ErrICSFailed)
	}
	return nil
}

func newJobName() string {
	id := rand.String(idLength)
	return "kn-event-sender-" + id
}
