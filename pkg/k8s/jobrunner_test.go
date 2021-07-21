package k8s_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestJobRunnerRun(t *testing.T) {
	clients := &tests.FakeClients{TB: t, Objects: make([]runtime.Object, 0)}
	runner := k8s.CreateJobRunner(clients)
	job := examplePiJob()
	jobs := clients.Typed().BatchV1().Jobs(job.Namespace)
	ctx := clients.Context()
	watcher, err := jobs.Watch(ctx, metav1.ListOptions{})
	assert.NoError(t, err)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.NoError(t, runner.Run(&job))
	}()
	<-watcher.ResultChan()
	watcher.Stop()
	sucJob := jobSuccess(job)
	_, err = jobs.Update(ctx, &sucJob, metav1.UpdateOptions{})
	assert.NoError(t, err)
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
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
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
