/*
 Copyright 2024 The Knative Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/pkg/ptr"
)

type jobGatherer struct {
	kube Clients
}

func (g jobGatherer) gather(ctx context.Context, job *batchv1.Job) outlogging.Fields {
	fields := outlogging.Fields{}
	gatherImageInfo(job, fields)
	asJSONIntoFields(job, "job", fields)
	// empty status -> the job isn't started yet
	if !reflect.DeepEqual(job.Status, batchv1.JobStatus{}) {
		// collect for a job that has been executed only
		g.gatherInfoOfPodsForJob(ctx, job, fields)
	}
	return fields
}

func gatherImageInfo(job *batchv1.Job, fields outlogging.Fields) {
	switch len(job.Spec.Template.Spec.Containers) {
	case 0:
		// nothing
	case 1:
		fields["image"] = job.Spec.Template.Spec.Containers[0].Image
	default:
		imgs := make([]string, len(job.Spec.Template.Spec.Containers))
		for i, container := range job.Spec.Template.Spec.Containers {
			imgs[i] = container.Image
		}
		asJSONIntoFields(imgs, "images", fields)
	}
}

func (g jobGatherer) gatherInfoOfPodsForJob(
	ctx context.Context, job *batchv1.Job, fields outlogging.Fields,
) {
	podsClient := g.kube.Typed().CoreV1().Pods(job.GetNamespace())
	list, err := podsClient.List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + job.GetName()})
	if err != nil {
		fields["podlist-err"] = err.Error()
	} else {
		asJSONIntoFields(list, "pods", fields)
		collectPodLogs(ctx, list, podsClient, fields)
		eventsClient := g.kube.Typed().CoreV1().Events(job.GetNamespace())
		collectPodEvents(ctx, list, eventsClient, fields)
	}
}

func collectPodLogs(
	ctx context.Context, list *corev1.PodList, podsClient typedcorev1.PodInterface,
	fields outlogging.Fields,
) {
	logs := make(map[string]string, list.Size())
	for _, pod := range list.Items {
		podLogs, plerr := podsClient.GetLogs(pod.GetName(), &corev1.PodLogOptions{}).DoRaw(ctx)
		if plerr != nil {
			fields["podlogs-err"] = plerr.Error()
			break
		}
		logs[pod.GetName()] = string(podLogs)
	}
	asJSONIntoFields(logs, "logs", fields)
}

func collectPodEvents(
	ctx context.Context, list *corev1.PodList, eventsClient typedcorev1.EventInterface,
	fields outlogging.Fields,
) {
	events := make(map[string]*corev1.EventList, list.Size())
	for _, pod := range list.Items {
		selector := eventsClient.GetFieldSelector(
			ptr.String(pod.GetName()), ptr.String(pod.GetNamespace()),
			nil, ptr.String(string(pod.UID)),
		).String()
		eventList, err := eventsClient.List(ctx, metav1.ListOptions{
			FieldSelector: selector,
		})
		if err != nil {
			fields["events-err"] = err.Error()
			break
		}
		events[pod.GetName()] = eventList
	}
	asJSONIntoFields(events, "events", fields)
}

func asJSONIntoFields(obj any, label string, fields outlogging.Fields) {
	if bytes, err := json.Marshal(obj); err == nil {
		fields[label] = string(bytes)
	} else {
		fields[label] = fmt.Sprintf("%#v", obj)
		errLabel := label + "-json-marshal-err"
		fields[errLabel] = err.Error()
	}
}
