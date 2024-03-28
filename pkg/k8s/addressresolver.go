package k8s

import (
	"context"
	"fmt"
	"net/url"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
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
	ResolveAddress(ref *tracker.Reference, uri *apis.URL) (*url.URL, error)
}

// CreateAddressResolver will create ReferenceAddressResolver, or return an
// error.
func CreateAddressResolver(kube Clients) ReferenceAddressResolver {
	ctx := ctxWithDynamic(kube)
	return &addressResolver{
		kube: kube, ctx: addressable.WithDuck(ctx),
	}
}

type addressResolver struct {
	kube Clients
	ctx  context.Context
}

// ResolveAddress of a tracker.Reference with given uri (as apis.URL).
func (a *addressResolver) ResolveAddress(
	ref *tracker.Reference,
	uri *apis.URL,
) (*url.URL, error) {
	gvr := a.toGVR(ref)
	dest, err := a.toDestination(gvr, ref, uri)
	if err != nil {
		return nil, err
	}
	parent := toAccessor(ref)
	tr := tracker.New(noopCallback, controller.GetTrackerLease(a.ctx))
	r := resolver.NewURIResolverFromTracker(a.ctx, tr)
	u, err := r.URIFromDestinationV1(a.ctx, *dest, parent)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotAddressable, err)
	}
	resolved := u.URL()
	return resolved, nil
}

func (a *addressResolver) toDestination(
	gvr schema.GroupVersionResource,
	ref *tracker.Reference,
	uri *apis.URL,
) (*duckv1.Destination, error) {
	dest := &duckv1.Destination{
		Ref: &duckv1.KReference{
			Kind:       ref.Kind,
			Namespace:  ref.Namespace,
			Name:       ref.Name,
			APIVersion: ref.APIVersion,
		},
		URI: uri,
	}
	if ref.Selector != nil {
		list, err := a.kube.Dynamic().Resource(gvr).
			Namespace(ref.Namespace).List(a.ctx, metav1.ListOptions{
			LabelSelector: ref.Selector.String(),
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		count := len(list.Items)
		if count == 0 {
			return nil, ErrNotFound
		}
		if count > 1 {
			return nil, fmt.Errorf("%w: %d", ErrMoreThenOneFound, count)
		}
		dest.Ref.Name = list.Items[0].GetName()
	}
	return dest, nil
}

func (a *addressResolver) toGVR(ref *tracker.Reference) schema.GroupVersionResource {
	gvk := ref.GroupVersionKind()
	gvr := apis.KindToResource(gvk)
	return gvr
}

func toAccessor(ref *tracker.Reference) kmeta.Accessor {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": ref.APIVersion,
		"kind":       ref.Kind,
		"metadata": map[string]interface{}{
			"name":      ref.Name,
			"namespace": ref.Namespace,
		},
	}}
}

func ctxWithDynamic(kube Clients) context.Context {
	return context.WithValue(kube.Context(), dynamicclient.Key{}, kube.Dynamic())
}

func noopCallback(_ types.NamespacedName) {
}
