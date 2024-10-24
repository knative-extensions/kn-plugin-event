package k8s_test

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"knative.dev/client/pkg/output"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
	testlogging "knative.dev/pkg/logging/testing"
	"knative.dev/pkg/ptr"
)

func TestJobRunnerRun(t *testing.T) {
	clients := &tests.FakeClients{TB: t, Objects: make([]runtime.Object, 0)}
	runner := k8s.NewJobRunner(clients)
	job := examplePiJob()
	jobs := clients.Typed().BatchV1().Jobs(job.Namespace)
	ctx, cancel := context.WithTimeout(testlogging.TestContextWithLogger(t), time.Second)
	printer := output.NewTestPrinter()
	defer cancel()
	ctx = output.WithContext(ctx, printer)
	watcher, err := jobs.Watch(ctx, metav1.ListOptions{})
	require.NoError(t, err)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.NoError(t, runner.Run(ctx, &job))
		// sleep a little, so the TUI had a chance to print the output
		time.Sleep(time.Millisecond)
	}()
	ev := <-watcher.ResultChan()
	assert.Equal(t, watch.Added, ev.Type)
	jobc, ok := ev.Object.(*batchv1.Job)
	assert.True(t, ok)
	assert.Equal(t, jobc.Name, job.GetName())
	watcher.Stop()
	sucJob := jobSuccess(job)
	_, err = jobs.Update(ctx, &sucJob, metav1.UpdateOptions{})
	require.NoError(t, err)
	wg.Wait()
	time.Sleep(time.Millisecond)
	assert.Contains(t, printer.Outputs().Out.String(),
		"Sending event within the cluster Done")
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
	name := "test-" +
		strconv.FormatInt(rand.Int63(), 36) //nolint:gosec
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "demo",
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: ptr.Int32(3),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ActiveDeadlineSeconds: ptr.Int64(10),
					Containers: []corev1.Container{{
						Image:   "docker.io/library/perl",
						Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
					}},
				},
			},
		},
	}
}
