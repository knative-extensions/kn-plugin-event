package sender

import (
	"fmt"
	"regexp"

	"github.com/cardil/kn-event/internal/cli/ics"
	"github.com/cardil/kn-event/internal/event"
	"github.com/cardil/kn-event/internal/k8s"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
		return fmt.Errorf("%w: %v", k8s.ErrInvalidReference, err)
	}
	kevent, err := ics.Encode(ce)
	if err != nil {
		return fmt.Errorf("%w: %v", ics.ErrCouldntEncode, err)
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("kn-event-sender-%s", ce.ID()),
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
						Image: imageFor("kn-event-sender"),
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
	return i.jobRunner.Run(job)
}

func imageFor(artifact string) string {
	basename := ics.ContainerBasename
	r := regexp.MustCompile(".+[A-Za-z0-9]$")
	if r.MatchString(basename) {
		basename += "/"
	}
	return basename + artifact
}
