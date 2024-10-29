package k8s_test

import (
	"context"
	"fmt"
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
	typedbatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"knative.dev/client/pkg/output"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
	testlogging "knative.dev/pkg/logging/testing"
	"knative.dev/pkg/ptr"
)

func TestJobRunnerRun(t *testing.T) {
	t.Parallel()
	tcs := []testJobRunnerRunTestCase{{
		"successful",
		1,
		nil,
	}, {
		"failure",
		3,
		k8s.ErrJobFailed,
	}}
	for _, tc := range tcs {
		t.Run(tc.name, testJobRunnerRun(tc))
	}
}

type testJobRunnerRunTestCase struct {
	name  string
	fails int
	err   error
}

func testJobRunnerRun(tc testJobRunnerRunTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		execs := tc.fails + 1
		job := examplePiJob()
		objs := make([]runtime.Object, 0, execs)
		for i := range execs {
			objs = append(objs, podForJob(i, job))
		}
		clients := &tests.FakeClients{TB: t, Objects: objs}
		runner := k8s.NewJobRunner(clients)
		jobs := clients.Typed().BatchV1().Jobs(job.Namespace)
		ctx, cancel := context.WithTimeout(testlogging.TestContextWithLogger(t), time.Second)
		printer := output.NewTestPrinter()
		defer cancel()
		ctx = output.WithContext(ctx, printer)
		watcher, err := jobs.Watch(ctx, metav1.ListOptions{})
		require.NoError(t, err)
		wg := sync.WaitGroup{}
		wg.Add(1)
		errs := make(chan error)
		go func() {
			defer wg.Done()
			if err := runner.Run(ctx, &job); err != nil {
				errs <- err
			}
			close(errs)
			// sleep a little, so the TUI had a chance to print the output
			time.Sleep(time.Millisecond)
		}()
		ev := <-watcher.ResultChan()
		assert.Equal(t, watch.Added, ev.Type)
		jobc, ok := ev.Object.(*batchv1.Job)
		assert.True(t, ok)
		assert.Equal(t, jobc.Name, job.GetName())
		watcher.Stop()

		completeJob(ctx, t, tc, job, jobs)
		err = <-errs
		wg.Wait()

		assert.Contains(t, printer.Outputs().Out.String(),
			"Sending event within the cluster")
		require.ErrorIs(t, err, tc.err)
		if tc.err == nil {
			assert.Contains(t, printer.Outputs().Out.String(),
				"Done")
		}
	}
}

func completeJob(
	ctx context.Context, t *testing.T, tc testJobRunnerRunTestCase,
	job batchv1.Job, jobs typedbatchv1.JobInterface,
) {
	t.Helper()
	for range tc.fails {
		job = addFailure(job)
		_, err := jobs.Update(ctx, &job, metav1.UpdateOptions{})
		require.NoError(t, err)
	}
	if tc.fails < int(*job.Spec.BackoffLimit) {
		job = completeSuccessfully(job)
		_, err := jobs.Update(ctx, &job, metav1.UpdateOptions{})
		require.NoError(t, err)
	}
}

func podForJob(i int, job batchv1.Job) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprint("run-", i),
			Namespace: job.GetNamespace(),
			Labels: map[string]string{
				"job-name": job.GetName(),
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Image: "example.org/foo/bar",
			}},
		},
	}
}

func addFailure(job batchv1.Job) batchv1.Job {
	job.Status.Failed++
	return job
}

func completeSuccessfully(job batchv1.Job) batchv1.Job {
	now := metav1.Now()
	job.Status.Succeeded++
	job.Status.Active = 0
	job.Status.Failed = 0
	job.Status.CompletionTime = &now
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
