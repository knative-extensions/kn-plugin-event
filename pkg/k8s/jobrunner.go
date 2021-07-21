package k8s

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// JobRunner will launch a Job and monitor it for completion.
type JobRunner interface {
	Run(*batchv1.Job) error
}

// CreateJobRunner will create a JobRunner, or return an error.
func CreateJobRunner(kube Clients) JobRunner {
	return &jobRunner{
		kube: kube,
	}
}

type jobRunner struct {
	kube Clients
}

func (j *jobRunner) Run(job *batchv1.Job) error {
	ctx := j.kube.Context()
	jobs := j.kube.Typed().BatchV1().Jobs(job.Namespace)
	_, err := jobs.Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
	err = j.watchJob(job, func(job *batchv1.Job) (bool, error) {
		if job.Status.CompletionTime == nil && job.Status.Failed == 0 {
			return false, nil
		}
		// We should be done if we reach here.
		if job.Status.Succeeded < 1 {
			return false, fmt.Errorf("%w: %s", ErrICSenderJobFailed,
				"expected to have successful job")
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
	return j.deleteJob(job)
}

func (j *jobRunner) deleteJob(job *batchv1.Job) error {
	ctx := j.kube.Context()
	jobs := j.kube.Typed().BatchV1().Jobs(job.GetNamespace())
	err := jobs.Delete(ctx, job.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
	pods := j.kube.Typed().CoreV1().Pods(job.GetNamespace())
	err = pods.DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job.GetName()),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
	return nil
}

func (j *jobRunner) watchJob(meta metav1.Object, changeFn func(job *batchv1.Job) (bool, error)) error {
	ctx := j.kube.Context()
	jobs := j.kube.Typed().BatchV1().Jobs(meta.GetNamespace())
	watcher, err := jobs.Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", meta.GetName()),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
	defer watcher.Stop()
	for result := range watcher.ResultChan() {
		if result.Type == watch.Added || result.Type == watch.Modified {
			job, ok := result.Object.(*batchv1.Job)
			if !ok {
				return fmt.Errorf("%w: %s: %T", ErrICSenderJobFailed,
					"expected to watch batchv1.Job, got", result.Object)
			}
			var brk bool
			brk, err = changeFn(job)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
			}
			if brk {
				watcher.Stop()
			}
		}
	}
	return nil
}
