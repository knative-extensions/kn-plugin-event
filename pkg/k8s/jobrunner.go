package k8s

import (
	"fmt"
	"sync"

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

type task struct {
	errs  chan<- error
	ready chan<- bool
	wg    *sync.WaitGroup
}

func (j *jobRunner) Run(job *batchv1.Job) error {
	ready := make(chan bool)
	errs := make(chan error)
	tsk := task{
		errs, ready, &sync.WaitGroup{},
	}
	tasks := []func(*batchv1.Job, task){
		// wait is started first,  making sure to capture success, even the ultra-fast one.
		j.waitForSuccess,
		j.createJob,
	}
	tsk.wg.Add(len(tasks))
	// run all tasks in parallel
	for _, fn := range tasks {
		go fn(job, tsk)
		<-ready
	}
	go waitAndClose(tsk)
	// return the first error
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return j.deleteJob(job)
}

func (j *jobRunner) createJob(job *batchv1.Job, tsk task) {
	defer tsk.wg.Done()
	tsk.ready <- true
	ctx := j.kube.Context()
	jobs := j.kube.Typed().BatchV1().Jobs(job.Namespace)
	_, err := jobs.Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		tsk.errs <- fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
}

func (j *jobRunner) waitForSuccess(job *batchv1.Job, tsk task) {
	defer tsk.wg.Done()
	err := j.watchJob(job, tsk, func(job *batchv1.Job) (bool, error) {
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
		tsk.errs <- fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
}

func waitAndClose(tsk task) {
	tsk.wg.Wait()
	close(tsk.errs)
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

func (j *jobRunner) watchJob(obj metav1.Object, tsk task, changeFn func(job *batchv1.Job) (bool, error)) error {
	ctx := j.kube.Context()
	jobs := j.kube.Typed().BatchV1().Jobs(obj.GetNamespace())
	watcher, err := jobs.Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", obj.GetName()),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrICSenderJobFailed, err)
	}
	defer watcher.Stop()
	resultCh := watcher.ResultChan()
	tsk.ready <- true
	for result := range resultCh {
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
				return nil
			}
		}
	}
	return nil
}
