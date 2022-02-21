package k8s_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestJobRunnerRun(t *testing.T) {
	clients := &tests.FakeClients{TB: t, Objects: make([]runtime.Object, 0)}
	runner := k8s.CreateJobRunner(clients)
	job := examplePiJob()
	jobs := clients.Typed().BatchV1().Jobs(job.Namespace)
	ctx := context.TODO()
	watcher, err := jobs.Watch(ctx, metav1.ListOptions{})
	assert.NilError(t, err)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.NilError(t, runner.Run(ctx, &job))
	}()
	ev := <-watcher.ResultChan()
	assert.Equal(t, ev.Type, watch.Added)
	assert.Equal(t, ev.Object.(*batchv1.Job).Name, job.GetName()) // nolint:forcetypeassert
	watcher.Stop()
	sucJob := jobSuccess(job)
	_, err = jobs.Update(ctx, &sucJob, metav1.UpdateOptions{})
	assert.NilError(t, err)
	wg.Wait()
}

func jobSuccess(job batchv1.Job) batchv1.Job {
	now := metav1.Now()
	job.Status.Succeeded = 1
	job.Status.Active = 0
	job.Status.Failed = 0
	job.Status.CompletionTime = &now
	job.Status.StartTime = &now
	job.Status.Conditions = []batchv1.JobCondition{{
		Type:               batchv1.JobComplete,
		Status:             corev1.ConditionTrue,
		LastProbeTime:      now,
		LastTransitionTime: now,
		Reason:             "done",
		Message:            "success",
	}}
	return job
}

func examplePiJob() batchv1.Job {
	name := fmt.Sprintf("test-%s",
		strconv.FormatInt(rand.Int63(), 36)) //nolint:gosec
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "demo",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "docker.io/library/perl",
						Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
					}},
				},
			},
		},
	}
}
