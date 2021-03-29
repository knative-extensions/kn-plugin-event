package internal

import "github.com/wavesoftware/go-ensure"

// gitVersion returns a git version string.
func gitVersion() string {
	if gitVerCache == nil {
		ver, err := git("describe", "--always", "--tags", "--dirty")
		ensure.NoError(err)
		gitVerCache = &ver
	}
	return *gitVerCache
}
