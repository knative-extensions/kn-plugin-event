package git

import (
	"context"

	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/go-magetasks/config"
)

var cacheKey = struct{}{}

// Version returns a git version string.
func Version() string {
	if version, ok := fromContext(); ok {
		return version
	}
	version, err := sh.Output("git", "describe",
		"--always", "--tags", "--dirty")
	ensure.NoError(err)
	saveInContext(version)
	return version
}

func saveInContext(version string) {
	config.WithContext(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, cacheKey, version)
	})
}

func fromContext() (string, bool) {
	ver, ok := config.Actual().Context.Value(cacheKey).(string)
	return ver, ok
}
