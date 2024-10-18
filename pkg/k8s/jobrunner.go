package k8s

import (
	"context"
	"fmt"
	"sync"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/pkg/errors"
	"sigs.k8s.io/yaml"
)

// JobRunner will launch a Job and monitor it for completion.
type JobRunner interface {
	Run(ctx context.Context, job *batchv1.Job) error
}

// NewJobRunner will create a JobRunner, or return an error.
func NewJobRunner(kube Clients) JobRunner {
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

func (j *jobRunner) Run(ctx context.Context, job *batchv1.Job) error {
	ready := make(chan bool)
	errs := make(chan error)
	tsk := task{
		errs, ready, &sync.WaitGroup{},
	}
	logdumpJob(ctx, job)

	tasks := []func(context.Context, *batchv1.Job, task){
		// wait is started first, making sure to capture success, even the ultra-fast one.
		j.waitForSuccess,
		j.createJob,
	}
	tsk.wg.Add(len(tasks))
	// run all tasks in parallel
	for _, fn := range tasks {
		go fn(ctx, job, tsk)
		<-ready
	}
	go waitAndClose(tsk)
	// return the first error
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return j.deleteJob(ctx, job)
}

func (j *jobRunner) createJob(ctx context.Context, job *batchv1.Job, tsk task) {
	defer tsk.wg.Done()
	tsk.ready <- true
	jobs := j.kube.Typed().BatchV1().Jobs(job.Namespace)
	_, err := jobs.Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		tsk.errs <- errors.Wrap(err, ErrJobFailed)
	}
}

func (j *jobRunner) waitForSuccess(ctx context.Context, job *batchv1.Job, tsk task) {
	defer tsk.wg.Done()
	err := j.watchJob(ctx, job, tsk, func(job *batchv1.Job) (bool, error) {
		if job.Status.Succeeded >= 1 {
			return true, nil
		}
		limit := int32(0)
		if job.Spec.BackoffLimit != nil {
			limit = *job.Spec.BackoffLimit
		}
		if job.Status.Failed >= limit {
			logdumpJob(ctx, job)
			return false, fmt.Errorf(
				"%w \"%s\" %d times, exceeding the limit of %d",
				ErrJobFailed, job.GetName(), job.Status.Failed, limit)
		}
		return false, nil
	})
	if err != nil {
		tsk.errs <- errors.Wrap(err, ErrJobFailed)
	}
}

func waitAndClose(tsk task) {
	tsk.wg.Wait()
	close(tsk.errs)
}

func (j *jobRunner) deleteJob(ctx context.Context, job *batchv1.Job) error {
	jobs := j.kube.Typed().BatchV1().Jobs(job.GetNamespace())
	policy := metav1.DeletePropagationBackground
	err := jobs.Delete(ctx, job.GetName(), metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		return errors.Wrap(err, ErrJobFailed)
	}
	pods := j.kube.Typed().CoreV1().Pods(job.GetNamespace())
	err = pods.DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "job-name=" + job.GetName(),
	})
	if err != nil {
		return errors.Wrap(err, ErrJobFailed)
	}
	return nil
}

func (j *jobRunner) watchJob(
	ctx context.Context,
	obj metav1.Object,
	tsk task,
	changeFn func(job *batchv1.Job) (bool, error),
) error {
	jobs := j.kube.Typed().BatchV1().Jobs(obj.GetNamespace())
	watcher, err := jobs.Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + obj.GetName(),
	})
	if err != nil {
		return errors.Wrap(err, ErrJobFailed)
	}
	defer watcher.Stop()
	resultCh := watcher.ResultChan()
	tsk.ready <- true
	for result := range resultCh {
		if result.Type == watch.Added || result.Type == watch.Modified {
			job, ok := result.Object.(*batchv1.Job)
			if !ok {
				return fmt.Errorf("%w: %s: %T", ErrJobFailed,
					"expected to watch batchv1.Job, got", result.Object)
			}
			var brk bool
			brk, err = changeFn(job)
			if err != nil {
				return errors.Wrap(err, ErrJobFailed)
			}
			if brk {
				return nil
			}
		}
	}
	return nil
}

func logdumpJob(ctx context.Context, job *batchv1.Job) {
	log := outlogging.LoggerFrom(ctx)
	image := "<unknown>"
	if len(job.Spec.Template.Spec.Containers) == 1 {
		image = job.Spec.Template.Spec.Containers[0].Image
	}
	if jobBytes, yerr := yaml.Marshal(job); yerr == nil {
		log.WithFields(outlogging.Fields{"job": string(jobBytes)}).
			Debug("Sender job image: ", image)
	} else {
		log.WithFields(outlogging.Fields{"err": yerr.Error()}).
			Debug("Sender job image: ", image)
	}
}
