package magetasks

import (
	"os"

	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Clean will clean project files.
func Clean() {
	t := tasks.Start("🚿", "Cleaning")
	err := os.RemoveAll(internal.BuildDir())
	errs := make([]error, 0, 1)
	errs = append(errs, err)
	for _, task := range config.CleaningTasks {
		p := t.Part(task.Name)
		p.Starting().Done(task.Task())
	}
	t.End(errs...)
}
