package magetasks

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/internal"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Test will execute regular unit tests.
func Test() {
	mg.Deps(Check, internal.EnsureBuildDir)
	t := tasks.StartMultiline("✅", "Testing")
	cmd := "richgo"
	if color.NoColor {
		cmd = "go"
	}
	args := []string{
		"test", "-v", "-covermode=count",
		fmt.Sprintf("-coverprofile=%s/coverage.out", internal.BuildDir()),
	}
	args = internal.AppendLdflags(args, t)
	args = append(args, "./...")
	err := sh.RunV(cmd, args...)
	t.End(err)
}
