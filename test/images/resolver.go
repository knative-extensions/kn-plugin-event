package images

import (
	"github.com/google/go-containerregistry/pkg/name"
	"k8s.io/apimachinery/pkg/util/json"
	"knative.dev/reconciler-test/pkg/environment"
)

// Resolver will resolve given KO package paths into real OCI images references.
// This interface probably should be moved into reconciler-test framework. See:
// https://github.com/knative-sandbox/reconciler-test/issues/303
type Resolver interface {
	// Resolve will resolve given KO package path into real OCI image reference.
	Resolve(kopath string) (name.Reference, error)
	// Applicable will tell that given resolver is applicable to current runtime
	// environment, or not.
	Applicable() bool
}

// TestingT a subset of testing.T.
type TestingT interface {
	Logf(fmt string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

// ResolveImages will try to resolve the images, using given resolver(s).
func ResolveImages(t TestingT, packages []string) {
	for _, resolver := range Resolvers {
		if resolver.Applicable() {
			resolveImagesWithResolver(t, packages, resolver)
			return
		}
	}
	if len(Resolvers) > 0 {
		t.Fatalf("Couldn't resolve images with registered resolvers: %+q", Resolvers)
	}
}

func resolveImagesWithResolver(t TestingT, packages []string, resolver Resolver) {
	resolved := make(map[string]string)
	for _, pack := range packages {
		kopath := "ko://" + pack
		image, err := resolver.Resolve(kopath)
		if err != nil {
			t.Fatal(err)
		}
		resolved[kopath] = image.String()
	}
	repr, err := json.Marshal(resolved)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Images resolved to: %s", string(repr))
	environment.WithImages(resolved)
}
