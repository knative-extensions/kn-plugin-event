package tasks

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/wavesoftware/go-ensure"
)

var (
	red    = color.New(color.FgHiRed).Add(color.Bold).SprintFunc()
	green  = color.New(color.FgHiGreen).Add(color.Bold).SprintFunc()
	yellow = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
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
	Done(err error)
}

// Skip print a warning message.
func (p *Part) Skip(reason string) {
	if p.t.multiline {
		msg := fmt.Sprintf(" %s › %s is skipped due to %s\n", p.t.icon, p.name, reason)
		fmt.Print(mageTag() + yellow(msg))
	}
}

type partProcessing struct {
	p *Part
}

// Done is reporting a completeness of part processing.
func (p *partProcessing) Done(err error) {
	if err != nil {
		if p.p.t.multiline {
			msg := fmt.Sprintf(" %s %s have failed: %v\n", p.p.t.icon, p.p.name, err)
			fmt.Print(mageTag() + red(msg))
		}
		p.p.t.raised = append(p.p.t.raised, err)
	}
}

// Starting starts a part processing.
func (p *Part) Starting() PartProcessing {
	if p.t.multiline {
		msg := fmt.Sprintf(" %s › %s\n", p.t.icon, p.name)
		fmt.Print(mageTag() + msg)
	}
	return &partProcessing{p: p}
}

// Start will start a single line task.
func Start(icon, action string) *Task {
	t := &Task{
		icon:      icon,
		action:    action,
		multiline: false,
	}
	t.start()
	return t
}

// StartMultiline will start a multi line task.
func StartMultiline(icon, action string) *Task {
	t := &Task{
		icon:      icon,
		action:    action,
		multiline: true,
	}
	t.start()
	return t
}

func (t *Task) start() {
	if t.multiline {
		fmt.Printf("%s %s %s\n", mageTag(), t.icon, t.action)
	} else {
		fmt.Printf("%s %s %s... ", mageTag(), t.icon, t.action)
	}
}

// End will report task completion, either successful or failures.
func (t *Task) End(errs ...error) {
	var msg string
	sum := make([]error, 0, len(errs)+len(t.raised))
	sum = append(sum, t.raised...)
	sum = append(sum, errs...)
	merr := multierror.Append(nil, sum...)
	err := merr.ErrorOrNil()
	if err != nil {
		msg = erroneousMsg(t)
	} else {
		msg = successfulMsg(t)
	}

	fmt.Print(msg)
	ensure.NoError(err)
}

// Part create a part that can be reported further.
func (t *Task) Part(part string) *Part {
	return &Part{
		name: part,
		t:    t,
	}
}

func erroneousMsg(t *Task) string {
	if t.multiline {
		return mageTag() + red(fmt.Sprintf(" %s have failed!\n", t.action))
	}
	return red(fmt.Sprintln("failed!"))
}

func successfulMsg(t *Task) string {
	if t.multiline {
		return mageTag() + green(fmt.Sprintf(" %s was successful.\n", t.action))
	}
	return green(fmt.Sprintln("done."))
}
