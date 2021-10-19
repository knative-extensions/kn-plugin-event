package tasks

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/output"
	"github.com/wavesoftware/go-magetasks/pkg/output/color"
)

// Task represents a mage task enriched with icon and description that might
// be multiline.
type Task struct {
	icon      string
	action    string
	multiline bool
	raised    []error
}

// Part represents a part of a bigger Task.
type Part struct {
	name string
	t    *Task
}

// PartProcessing is an interface that is used to process long running part of
// tasks.
type PartProcessing interface {
	config.Notifier
	Done(err error)
}

// Skip print a warning message.
func (p *Part) Skip(reason string) {
	if p.t.multiline {
		msg := fmt.Sprintf("%s › %s is skipped due to %s", p.t.icon, p.name, reason)
		output.Println(color.Yellow(msg))
	}
}

type partProcessing struct {
	p *Part
}

func (pp *partProcessing) Notify(status string) {
	if pp.p.t.multiline {
		output.Printlnf("%s › %s › %s",
			pp.p.t.icon, pp.p.name, status)
	}
}

// Done is reporting a completeness of part processing.
func (pp *partProcessing) Done(err error) {
	if err != nil {
		if pp.p.t.multiline {
			msg := fmt.Sprintf("%s %s have failed: %v", pp.p.t.icon, pp.p.name, err)
			output.Println(color.Red(msg))
		}
		pp.p.t.raised = append(pp.p.t.raised, err)
	}
}

// Starting starts a part processing.
func (p *Part) Starting() PartProcessing {
	if p.t.multiline {
		msg := fmt.Sprintf("%s › %s", p.t.icon, p.name)
		output.Println(msg)
	}
	return &partProcessing{p: p}
}

// Start will start a single line task.
func Start(icon, action string, multiline bool) *Task {
	t := &Task{
		icon:      icon,
		action:    action,
		multiline: multiline,
	}
	t.start()
	return t
}

func (t *Task) start() {
	if t.multiline {
		output.Printlnf("%s %s", t.icon, t.action)
	} else {
		output.PrintPending(t.icon, " ", t.action, "... ")
	}
}

// End will report task completion, either successful or failures.
func (t *Task) End(errs ...error) {
	sum := make([]error, 0, len(errs)+len(t.raised))
	sum = append(sum, t.raised...)
	sum = append(sum, errs...)
	merr := multierror.Append(nil, sum...)
	err := merr.ErrorOrNil()
	if err != nil {
		erroneousMsg(t)
	} else {
		successfulMsg(t)
	}

	ensure.NoError(err)
}

// Part create a part that can be reported further.
func (t *Task) Part(part string) *Part {
	return &Part{
		name: part,
		t:    t,
	}
}

func erroneousMsg(t *Task) {
	if t.multiline {
		output.Println(color.Red(fmt.Sprintf("%s %s have failed!", t.icon, t.action)))
	} else {
		output.PrintEnd(color.Red("failed!"))
	}
}

func successfulMsg(t *Task) {
	if t.multiline {
		output.Println(color.Green(fmt.Sprintf("%s %s was successful.", t.icon, t.action)))
	} else {
		output.PrintEnd(color.Green("done."))
	}
}
