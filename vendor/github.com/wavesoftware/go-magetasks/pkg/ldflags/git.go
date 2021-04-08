package ldflags

import (
	"github.com/wavesoftware/go-magetasks/internal"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// AppendGitVersion will set version variable with git describe info.
func AppendGitVersion(args []string, t *tasks.Task) []string {
	return internal.AppendLdflags(args, t)
}
