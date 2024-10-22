package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/pkg/errors"
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
	j.logdumpJob(ctx, "Job to be executed", job)

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
			j.logdumpJob(ctx, "Successful job", job)
			return true, nil
		}
		limit := int32(0)
		if job.Spec.BackoffLimit != nil {
			limit = *job.Spec.BackoffLimit
		}
		if job.Status.Failed >= limit {
			j.logdumpJob(ctx, "Failed job", job)
			return false, fmt.Errorf(
				"%w %d times, exceeding the limit (job name: \"%s\")",
				ErrJobFailed, job.Status.Failed, job.GetName())
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

func (j *jobRunner) logdumpJob(ctx context.Context, label string, job *batchv1.Job) {
	log := outlogging.LoggerFrom(ctx)
	image := "<unknown>"
	if len(job.Spec.Template.Spec.Containers) == 1 {
		image = job.Spec.Template.Spec.Containers[0].Image
	}
	fields := outlogging.Fields{
		"image": image,
	}
	marshalIntoFields(job, "job", "jmerr", fields)
	// empty status -> the job isn't started yet
	if !reflect.DeepEqual(job.Status, batchv1.JobStatus{}) {
		// collect pods for a job that has been executed
		j.marshalPodsOfJob(ctx, job, fields)
	}
	log.WithFields(fields).Debug(label)
}

func (j *jobRunner) marshalPodsOfJob(ctx context.Context, job *batchv1.Job, fields outlogging.Fields) {
	pods := j.kube.Typed().CoreV1().Pods(job.GetNamespace())
	list, err := pods.List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + job.GetName()})
	if err != nil {
		fields["perr"] = err.Error()
	} else {
		marshalIntoFields(list, "pods", "pmerr", fields)
		logs := make(map[string]string, list.Size())
		for _, pod := range list.Items {
			podLogs, plerr := pods.GetLogs(pod.GetName(), &corev1.PodLogOptions{}).DoRaw(ctx)
			if plerr != nil {
				fields["plerr"] = plerr.Error()
				break
			} else {
				logs[pod.GetName()] = string(podLogs)
			}
		}
		marshalIntoFields(logs, "logs", "lmerr", fields)
	}
}

func marshalIntoFields(obj any, label, errLabel string, fields outlogging.Fields) {
	if bytes, err := json.Marshal(obj); err == nil {
		fields[label] = string(bytes)
	} else {
		fields[label] = fmt.Sprintf("%#v", obj)
		fields[errLabel] = err.Error()
	}
}
