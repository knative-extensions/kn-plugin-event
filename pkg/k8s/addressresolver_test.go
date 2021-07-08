package k8s_test

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	clienttest "knative.dev/client/lib/test"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
	netpkg "knative.dev/networking/pkg"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/tracker"
	"knative.dev/serving/pkg/apis/serving"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestResolveAddress(t *testing.T) {
	ns := clienttest.NextNamespace()
	testCases := []resolveAddressTestCase{
		k8sService(ns),
		knServiceClusterLocal(ns),
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// FIXME: fails with: no kind is registered for the type v1.Service in scheme "pkg/tests/fakeclients.go:33"
			t.Skip("FIXME: fails with: no kind is registered for the type v1.Service in scheme \"pkg/tests/fakeclients.go:33\"")
			performResolveAddressTest(t, tc, func(tr resolveAddressTestResources) (k8s.Clients, func(t *testing.T)) {
				objects := make([]runtime.Object, 0, len(tr.k8sServices))
				for j := range tr.k8sServices {
					//goland:noinspection GoShadowedVar
					service := tr.k8sServices[j]
					objects = append(objects, &service)
				}
				for j := range tr.knServices {
					service := tr.knServices[j]
					objects = append(objects, &service)
				}
				return &tests.FakeClients{Objects: objects}, func(t *testing.T) {
					t.Helper() // do nothing
				}
			})
		})
	}
}

func performResolveAddressTest(
	t *testing.T,
	tc resolveAddressTestCase,
	clientsFn func(resolveAddressTestResources) (k8s.Clients, func(t *testing.T)),
) {
	clients, cleanup := clientsFn(tc.resolveAddressTestResources)
	defer cleanup(t)
	resolver := k8s.CreateAddressResolver(clients)
	uri, err := apis.ParseURL("/")
	assert.NilError(t, err)
	u, err := resolver.ResolveAddress(tc.ref, uri)
	if tc.err != nil {
		assert.ErrorType(t, err, tc.err)
	} else {
		assert.Equal(t, err, nil)
	}
	assert.Equal(t, tc.wantURL, u.String())
}

func k8sService(namespace string) resolveAddressTestCase {
	kind := "Service"
	apiVersion := "v1"
	return resolveAddressTestCase{
		name:    "k8s service",
		wantURL: fmt.Sprintf("http://k8s-hello.%s.svc/", namespace),
		err:     nil,
		ref: &tracker.Reference{
			APIVersion: apiVersion,
			Kind:       kind,
			Namespace:  namespace,
			Name:       "k8s-hello",
			Selector:   nil,
		},
		resolveAddressTestResources: resolveAddressTestResources{
			k8sServices: []corev1.Service{{
				TypeMeta: metav1.TypeMeta{
					Kind:       kind,
					APIVersion: apiVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "k8s-hello",
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Name: "http",
						Port: 80,
						TargetPort: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: 8080,
						},
					}},
					Type: corev1.ServiceTypeClusterIP,
				},
			}},
		},
	}
}

func knServiceClusterLocal(namespace string) resolveAddressTestCase {
	kind := "Service"
	apiVersion := "serving.knative.dev/v1"
	return resolveAddressTestCase{
		name:    "kn cluster local service",
		wantURL: fmt.Sprintf("http://kn-hello.%s.svc.cluster.local/", namespace),
		err:     nil,
		ref: &tracker.Reference{
			APIVersion: apiVersion,
			Kind:       kind,
			Namespace:  namespace,
			Name:       "kn-hello",
			Selector:   nil,
		},
		resolveAddressTestResources: resolveAddressTestResources{
			knServices: []servingv1.Service{{
				TypeMeta: metav1.TypeMeta{
					Kind:       kind,
					APIVersion: apiVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kn-hello",
					Namespace: namespace,
					Labels: map[string]string{
						netpkg.VisibilityLabelKey: serving.VisibilityClusterLocal,
					},
				},
				Spec: servingv1.ServiceSpec{
					ConfigurationSpec: servingv1.ConfigurationSpec{
						Template: servingv1.RevisionTemplateSpec{
							Spec: servingv1.RevisionSpec{
								PodSpec: corev1.PodSpec{
									Containers: []corev1.Container{{
										Image: "gcr.io/knative-samples/helloworld-go",
									}},
								},
							},
						},
					},
				},
			}},
		},
	}
}

type resolveAddressTestResources struct {
	k8sServices []corev1.Service
	knServices  []servingv1.Service
}

type resolveAddressTestCase struct {
	name    string
	wantURL string
	err     error
	ref     *tracker.Reference
	resolveAddressTestResources
}
