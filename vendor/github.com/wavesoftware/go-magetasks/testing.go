package magetasks

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/files"
	"github.com/wavesoftware/go-magetasks/pkg/ldflags"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Test will execute regular unit tests.
func Test() {
	mg.Deps(Check, files.EnsureBuildDir)
	t := tasks.Start("✅", "Testing", true)
	args := []string{
		"--format", "testname",
		"--jsonfile", fmt.Sprintf("%s/tests-output.json", files.BuildDir()),
		"--junitfile", fmt.Sprintf("%s/tests-results.junit.xml", files.BuildDir()),
	}
	if color.NoColor {
		args = append(args, "--no-color")
	}
	args = append(args, "--",
		"-covermode=atomic", fmt.Sprintf("-coverprofile=%s/coverage.out", files.BuildDir()),
		"-race",
		"-count=1",
		"-short",
	)
	args = append(appendBuildVariables(args), "./...")
	cmd := fmt.Sprintf("%s/tools/gotestsum", files.BuildDir())
	err := sh.RunV(cmd, args...)
	t.End(err)
}

func appendBuildVariables(args []string) []string {
	c := config.Actual()
	if c.Version != nil || len(c.BuildVariables) > 0 {
		builder := ldflags.NewBuilder()
		if c.Version != nil {
			builder.Add(c.Version.Path, c.Version.Resolver.Version)
		}
		for key, resolver := range c.BuildVariables {
			builder.Add(key, resolver)
		}
		args = builder.BuildOnto(args)
	}
	return args
}
