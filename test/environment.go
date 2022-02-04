package test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"gotest.tools/v3/icmd"
)

// ResolveKnEventCommand will look for the kn event executable. Executable can
// be overridden by setting the environment variable: KN_PLUGIN_EVENT_EXECUTABLE
// Base args can be overridden by setting the environment variable:
// KN_PLUGIN_EVENT_EXECUTABLE_ARGS.
//
// Set:
// ```
// KN_PLUGIN_EVENT_EXECUTABLE=/path/to/kn
// KN_PLUGIN_EVENT_EXECUTABLE_ARGS=event
// ```
// to test kn-plugin-event as embedded in kn CLI.
func ResolveKnEventCommand(t TestingT) Command {
	MaybeSkip(t)
	bin := fmt.Sprintf("kn-event-%s-%s", runtime.GOOS, runtime.GOARCH)
	c := Command{
		Executable: path.Join(rootDir(), "build", "_output", "bin", bin),
	}
	if val, ok := os.LookupEnv("KN_PLUGIN_EVENT_EXECUTABLE"); ok {
		c.Executable = val
	}
	if val, ok := os.LookupEnv("KN_PLUGIN_EVENT_EXECUTABLE_ARGS"); ok {
		c.Args = strings.Split(val, " ")
	}
	return c
}

// Command represents a binary command to be executed.
type Command struct {
	Executable string
	Args       []string
}

// ToIcmd converts to icmd.Cmd.
func (c Command) ToIcmd(args ...string) icmd.Cmd {
	args = append(c.Args, args...)
	return icmd.Command(c.Executable, args...)
}

func rootDir() string {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled
	return path.Dir(path.Dir(filename))
}

// MaybeSkip should be called for integration tests. It will skip test if
// testing's short flag is being set.
func MaybeSkip(t TestingT) {
	if testing.Short() {
		t.Skipf("Short flag is set. Skipping %s.", t.Name())
	}
}

// TestingT a subset of testing.T.
type TestingT interface {
	Name() string
	Skipf(format string, args ...interface{})
}
