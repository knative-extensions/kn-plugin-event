package images

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/json"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/environment"
)

// ErrCouldNotResolve is returned when image couldn't be resolved.
var ErrCouldNotResolve = errors.New("could not resolve image with registered resolvers")

// Resolver will resolve given KO package paths into real OCI images references.
// This interface probably should be moved into reconciler-test framework. See:
// https://github.com/knative-extensions/reconciler-test/issues/303
type Resolver interface {
	// Resolve will resolve given KO package path into real OCI image reference.
	Resolve(kopath string) (name.Reference, error)
	// Applicable will tell that given resolver is applicable to current runtime
	// environment, or not.
	Applicable() bool
}

// PackageResolver is a function that will return package name for given context.
type PackageResolver func(ctx context.Context) string

// ExplicitPackage will return given package.
func ExplicitPackage(pack string) PackageResolver {
	return func(context.Context) string {
		return pack
	}
}

// ResolveImages will try to resolve the images, using given resolver(s).
func ResolveImages(packages []PackageResolver) environment.EnvOpts {
	return func(ctx context.Context, _ environment.Environment) (context.Context, error) {
		for _, resolver := range Resolvers {
			if resolver.Applicable() {
				return resolveImagesWithResolver(ctx, resolver, packages)
			}
		}
		if len(Resolvers) > 0 {
			return nil, fmt.Errorf("%w: %+q", ErrCouldNotResolve, Resolvers)
		}
		return ctx, nil
	}
}

func resolveImagesWithResolver(
	ctx context.Context,
	resolver Resolver,
	packages []PackageResolver,
) (context.Context, error) {
	log := logging.FromContext(ctx)
	resolved := make(map[string]string)
	for _, packResolver := range packages {
		pack := packResolver(ctx)
		kopack := pack
		if !strings.HasPrefix(kopack, "ko://") {
			kopack = "ko://" + kopack
		}
		image, err := resolver.Resolve(kopack)
		if err != nil {
			log.Fatal(errors.WithStack(err))
		}
		resolved[kopack] = image.String()
	}
	repr, err := json.Marshal(resolved)
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}
	log.Infof("Images resolved to: %s", string(repr))
	opt := environment.WithImages(resolved)
	return opt(ctx, nil)
}
