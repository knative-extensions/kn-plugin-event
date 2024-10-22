package sender

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"knative.dev/kn-plugin-event/pkg/errors"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/ics"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

const (
	idLength          = 6
	retriesDefault    = 4
	retriesConfigKey  = "plugins.event.in-cluster-sender.retries"
	deadlineDefault   = 7 // seconds
	deadlineConfigKey = "plugins.event.in-cluster-sender.deadline"
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
	kevent, ierr := ics.Encode(ce)
	if ierr != nil {
		return errors.Wrap(ierr, ics.ErrCouldntEncode)
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
			BackoffLimit: retries(),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:         corev1.RestartPolicyNever,
					ActiveDeadlineSeconds: deadline(),
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

func deadline() *int64 {
	v := int64(deadlineDefault)
	if viper.IsSet(deadlineConfigKey) {
		v = viper.GetInt64(deadlineConfigKey)
	}
	return &v
}

func retries() *int32 {
	v := int32(retriesDefault)
	if viper.IsSet(retriesConfigKey) {
		v = viper.GetInt32(retriesConfigKey)
	}
	return &v
}

func newJobName() string {
	id := rand.String(idLength)
	return "kn-event-sender-" + id
}
