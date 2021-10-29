package knative

import (
	"github.com/wavesoftware/go-magetasks/pkg/cache"
	"github.com/wavesoftware/go-magetasks/pkg/environment"
	"github.com/wavesoftware/go-magetasks/pkg/git"
	"github.com/wavesoftware/go-magetasks/pkg/version"
)

// NewTestableVersionResolver creates an instance of version.Resolver that can
// be easily tested.
func NewTestableVersionResolver(
	repo git.Repository,
	env func() environment.Values,
) version.Resolver {
	return NewVersionResolver(
		WithGit(
			git.WithCache(cache.NoopCache{}),
			git.WithRepository(repo),
		),
		WithEnvironmental(environment.WithValuesSupplier(env)),
	)
}
