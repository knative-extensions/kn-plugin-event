package git

import (
	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/cache"
	"github.com/wavesoftware/go-magetasks/pkg/version"
)

type VersionResolverOption func(*VersionResolver)

func NewVersionResolver(options ...VersionResolverOption) VersionResolver {
	r := VersionResolver{}

	for _, option := range options {
		option(&r)
	}

	return r
}

func WithCache(cache cache.Cache) VersionResolverOption {
	return func(r *VersionResolver) {
		r.Cache = cache
	}
}

func WithIsLatestStrategy(strategy IsLatestStrategy) VersionResolverOption {
	return func(r *VersionResolver) {
		r.IsLatestStrategy = strategy
	}
}

func WithRepository(repository Repository) VersionResolverOption {
	return func(r *VersionResolver) {
		r.Repository = repository
	}
}

func WithRemote(remote Remote) VersionResolverOption {
	return func(r *VersionResolver) {
		r.Remote = &remote
	}
}

// VersionResolver implements version.Resolver for git SCM.
type VersionResolver struct {
	Cache cache.Cache
	IsLatestStrategy
	Repository
	*Remote
}

// Remote represents a remote repository name and address.
type Remote struct {
	Name string
	URL  string
}

// IsLatestStrategy is used to determine if current version is latest one.
type IsLatestStrategy func(version.Resolver) func(string) (bool, error)

type cacheKey struct {
	typee string
}

func (r VersionResolver) Version() string {
	ver, err := r.cache().Compute(cacheKey{"version"}, func() (interface{}, error) {
		return r.repository().Describe()
	})
	ensure.NoError(err)
	return ver.(string)
}

func (r VersionResolver) IsLatest(versionRange string) (bool, error) {
	return ResolveIsLatest(r, r, versionRange)
}

func (r VersionResolver) cache() cache.Cache {
	if r.Cache == nil {
		return config.Cache()
	}
	return r.Cache
}

func (r VersionResolver) repository() Repository {
	if r.Repository == nil {
		return installedGitBinaryRepo{r.remote()}
	}
	return r.Repository
}

func (r VersionResolver) remote() Remote {
	remote := Remote{Name: "origin"}
	if r.Remote != nil {
		remote = *r.Remote
	}
	return remote
}

func (r VersionResolver) resolveTags() []string {
	tt, err := r.cache().Compute(cacheKey{"tags"}, func() (interface{}, error) {
		return r.repository().Tags()
	})
	ensure.NoError(err)
	return tt.([]string)
}
