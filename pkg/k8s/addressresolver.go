package k8s

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/client/pkg/dynamic"
	"knative.dev/client/pkg/flags/sink"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/client/injection/ducks/duck/v1/addressable"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/resolver"
	"knative.dev/pkg/tracker"
)

// ReferenceAddressResolver will resolve the tracker.Reference to an url.URL, or
// return an error.
type ReferenceAddressResolver interface {
	ResolveAddress(ctx context.Context, ref *sink.Reference, relativeURI string) (*url.URL, error)
}

// NewAddressResolver will create ReferenceAddressResolver or return an
// error.
func NewAddressResolver(kube Clients) ReferenceAddressResolver {
	return &addressResolver{
		kube: kube,
	}
}

type addressResolver struct {
	kube Clients
}

// ResolveAddress of a tracker.Reference with given uri (as apis.URL).
func (a *addressResolver) ResolveAddress(
	ctx context.Context, ref *sink.Reference, relativeURI string,
) (*url.URL, error) {
	dest, err := ref.Resolve(ctx, a.knclients())
	if err != nil {
		return nil, err
	}
	if dest.URI != nil {
		return relativize(dest.URI, relativeURI), nil
	}
	parent := toAccessor(dest.Ref)
	ctx = context.WithValue(ctx, dynamicclient.Key{}, a.kube.Dynamic())
	ctx = addressable.WithDuck(ctx)
	tr := tracker.New(noopCallback, controller.GetTrackerLease(ctx))
	r := resolver.NewURIResolverFromTracker(ctx, tr)
	u, err := r.URIFromDestinationV1(ctx, *dest, parent)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotAddressable, err)
	}
	return relativize(u, relativeURI), nil
}

func relativize(uri *apis.URL, relativeURI string) *url.URL {
	if relativeURI == "" {
		return uri.URL()
	}
	u := uri.URL()
	u.Path = path.Clean(path.Join(u.Path, relativeURI))
	return u
}

func (a *addressResolver) knclients() dynamic.KnDynamicClient {
	return dynamic.NewKnDynamicClient(a.kube.Dynamic(), a.kube.Namespace())
}

func toAccessor(ref *duckv1.KReference) kmeta.Accessor {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": ref.APIVersion,
		"kind":       ref.Kind,
		"metadata": map[string]interface{}{
			"name":      ref.Name,
			"namespace": ref.Namespace,
		},
	}}
}

func noopCallback(_ types.NamespacedName) {
}
