/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1 "knative.dev/eventing/pkg/apis/flows/v1"
	scheme "knative.dev/eventing/pkg/client/clientset/versioned/scheme"
)

// ParallelsGetter has a method to return a ParallelInterface.
// A group's client should implement this interface.
type ParallelsGetter interface {
	Parallels(namespace string) ParallelInterface
}

// ParallelInterface has methods to work with Parallel resources.
type ParallelInterface interface {
	Create(ctx context.Context, parallel *v1.Parallel, opts metav1.CreateOptions) (*v1.Parallel, error)
	Update(ctx context.Context, parallel *v1.Parallel, opts metav1.UpdateOptions) (*v1.Parallel, error)
	UpdateStatus(ctx context.Context, parallel *v1.Parallel, opts metav1.UpdateOptions) (*v1.Parallel, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Parallel, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ParallelList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Parallel, err error)
	ParallelExpansion
}

// parallels implements ParallelInterface
type parallels struct {
	client rest.Interface
	ns     string
}

// newParallels returns a Parallels
func newParallels(c *FlowsV1Client, namespace string) *parallels {
	return &parallels{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the parallel, and returns the corresponding parallel object, and an error if there is any.
func (c *parallels) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Parallel, err error) {
	result = &v1.Parallel{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("parallels").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Parallels that match those selectors.
func (c *parallels) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ParallelList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ParallelList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("parallels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested parallels.
func (c *parallels) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("parallels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a parallel and creates it.  Returns the server's representation of the parallel, and an error, if there is any.
func (c *parallels) Create(ctx context.Context, parallel *v1.Parallel, opts metav1.CreateOptions) (result *v1.Parallel, err error) {
	result = &v1.Parallel{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("parallels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(parallel).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a parallel and updates it. Returns the server's representation of the parallel, and an error, if there is any.
func (c *parallels) Update(ctx context.Context, parallel *v1.Parallel, opts metav1.UpdateOptions) (result *v1.Parallel, err error) {
	result = &v1.Parallel{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("parallels").
		Name(parallel.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(parallel).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *parallels) UpdateStatus(ctx context.Context, parallel *v1.Parallel, opts metav1.UpdateOptions) (result *v1.Parallel, err error) {
	result = &v1.Parallel{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("parallels").
		Name(parallel.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(parallel).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the parallel and deletes it. Returns an error if one occurs.
func (c *parallels) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("parallels").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *parallels) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("parallels").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched parallel.
func (c *parallels) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Parallel, err error) {
	result = &v1.Parallel{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("parallels").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
