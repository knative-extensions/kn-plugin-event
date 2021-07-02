package k8s

import (
	"context"
	"fmt"
	"net/url"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/restmapper"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
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
	return &addressResolver{
		kube: kube, ctx: kube.Context(),
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
	gvr, err := a.toGVR(ref)
	if err != nil {
		return nil, err
	}
	dest, err := a.toDestination(gvr, ref, uri)
	if err != nil {
		return nil, err
	}
	un, err := a.kube.Dynamic().Resource(gvr).
		Namespace(ref.Namespace).Get(a.ctx, dest.Ref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	addr, err := a.toAddressable(un)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotAddressable, err)
	}
	u := addr.URL.ResolveReference(uri).URL()
	return u, nil
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
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
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

func (a *addressResolver) toGVR(ref *tracker.Reference) (schema.GroupVersionResource, error) {
	gvk := ref.GroupVersionKind()
	dc := a.kube.Typed().Discovery()
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(
		memory.NewMemCacheClient(dc),
	)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{},
			fmt.Errorf("%w: %v", ErrInvalidReference, err)
	}
	return mapping.Resource, nil
}

func (a *addressResolver) toAddressable(un *unstructured.Unstructured) (*duckv1.Addressable, error) {
	gvk := un.GroupVersionKind()
	if gvk.Version == "v1" && gvk.Kind == "Service" && gvk.Group == "" {
		return &duckv1.Addressable{
			URL: apis.HTTP(fmt.Sprintf("%s.%s.svc", un.GetName(), un.GetNamespace())),
		}, nil
	}
	addr := &duckv1.Addressable{}
	err := duck.VerifyType(un, addr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotAddressable, err)
	}
	return addr, nil
}
