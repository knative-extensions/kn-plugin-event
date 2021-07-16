package k8s_test

import (
	"fmt"
	"net/url"
	"strings"
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
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/tracker"
	"knative.dev/serving/pkg/apis/serving"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestResolveAddress(t *testing.T) {
	ns := clienttest.NextNamespace()
	resolveAddressTestCases(ns, func(tc resolveAddressTestCase) {
		t.Run(tc.name, func(t *testing.T) {
			performResolveAddressTest(t, tc, func(tr resolveAddressTestResources) (k8s.Clients, func(tb testing.TB)) {
				return fakeClients(t, tr), noCleanup
			})
		})
	})
}

func noCleanup(tb testing.TB) {
	tb.Helper() // do nothing
}

func performResolveAddressTest( //nolint:thelper
	// lint skipped for greater visibility of failure place
	tb testing.TB,
	tc resolveAddressTestCase,
	clientsFn func(resolveAddressTestResources) (k8s.Clients, func(tb testing.TB)),
) {
	uri, err := apis.ParseURL("/")
	assert.NilError(tb, err)
	clients, cleanup := clientsFn(tc.resolveAddressTestResources)
	defer cleanup(tb)
	resolver := k8s.CreateAddressResolver(clients)
	u, err := resolver.ResolveAddress(tc.ref, uri)
	if tc.err != nil {
		assert.ErrorType(tb, err, tc.err)
	} else {
		assert.Equal(tb, err, nil)
	}
	assert.Check(tb, tc.matches(u), u)
}

func resolveAddressTestCases(namespace string, casefn func(tc resolveAddressTestCase)) {
	tcs := []resolveAddressTestCase{
		k8sService(namespace),
		knService(namespace, true),
		knService(namespace, false),
	}
	for _, tc := range tcs {
		casefn(tc)
	}
}

func k8sService(namespace string) resolveAddressTestCase {
	kind := "Service"
	apiVersion := "v1"
	return resolveAddressTestCase{
		name: "k8s service",
		matches: func(u *url.URL) bool {
			return u.String() ==
				fmt.Sprintf("http://k8s-hello.%s.svc.cluster.local/", namespace)
		},
		err: nil,
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

func knService(namespace string, clusterLocal bool) resolveAddressTestCase {
	m := matcher{
		local:     clusterLocal,
		name:      "kn-hello",
		namespace: namespace,
	}
	labels := map[string]string{}
	if clusterLocal {
		m.name = fmt.Sprintf("%s-cl", m.name)
		labels[netpkg.VisibilityLabelKey] = serving.VisibilityClusterLocal
	}
	clusterLocalURL := apis.HTTP(fmt.Sprintf(
		"%s.%s.svc.cluster.local", m.name, namespace))
	clusterLocalURL.Path = "/"
	m.url = clusterLocalURL
	publicURL := apis.HTTP(fmt.Sprintf(
		"%s-%s.apps.cloud.example.org", m.name, namespace))
	publicURL.Path = "/"
	serviceURL := publicURL
	if clusterLocal {
		serviceURL = clusterLocalURL
	}
	kind := "Service"
	apiVersion := "serving.knative.dev/v1"
	return resolveAddressTestCase{
		name:    m.name,
		matches: m.matches,
		err:     nil,
		ref: &tracker.Reference{
			APIVersion: apiVersion,
			Kind:       kind,
			Namespace:  namespace,
			Name:       m.name,
			Selector:   nil,
		},
		resolveAddressTestResources: resolveAddressTestResources{
			knServices: []servingv1.Service{{
				TypeMeta: metav1.TypeMeta{
					Kind:       kind,
					APIVersion: apiVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      m.name,
					Namespace: namespace,
					Labels:    labels,
				},
				Spec: ksvcSpec("gcr.io/knative-samples/helloworld-go"),
				Status: servingv1.ServiceStatus{
					RouteStatusFields: servingv1.RouteStatusFields{
						URL: serviceURL,
						Address: &duckv1.Addressable{
							URL: clusterLocalURL,
						},
					},
				},
			}},
		},
	}
}

type matcher struct {
	url       *apis.URL
	local     bool
	name      string
	namespace string
}

func (m matcher) matches(u *url.URL) bool {
	if m.local {
		return u.String() == m.url.String()
	}
	return strings.Contains(u.Host, m.name) &&
		strings.Contains(u.Host, m.namespace) &&
		u.String() != m.url.String()
}

func ksvcSpec(image string) servingv1.ServiceSpec {
	return servingv1.ServiceSpec{
		ConfigurationSpec: servingv1.ConfigurationSpec{
			Template: servingv1.RevisionTemplateSpec{
				Spec: servingv1.RevisionSpec{
					PodSpec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Image: image,
						}},
					},
				},
			},
		},
	}
}

func fakeClients(tb testing.TB, tr resolveAddressTestResources) k8s.Clients {
	tb.Helper()
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
	return &tests.FakeClients{
		Objects: objects,
		TB:      tb,
	}
}

type resolveAddressTestResources struct {
	k8sServices []corev1.Service
	knServices  []servingv1.Service
}

type resolveAddressTestCase struct {
	name    string
	matches func(url *url.URL) bool
	err     error
	ref     *tracker.Reference
	resolveAddressTestResources
}
