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
	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
	netpkg "knative.dev/networking/pkg"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/tracker"
	"knative.dev/serving/pkg/apis/serving"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestResolveAddress(t *testing.T) {
	ns := clienttest.NextNamespace()
	resolveAddressTestCases(ns, func(tc resolveAddressTestCase) {
		t.Run(tc.name, func(t *testing.T) {
			performResolveAddressTest(t, tc, func() (k8s.Clients, func(tb testing.TB)) {
				return fakeClients(t, tc), noCleanup
			})
		})
	})
}

func noCleanup(tb testing.TB) {
	tb.Helper() // do nothing
}

// performResolveAddressTest thelper lint skipped for greater visibility of
// failure location.
func performResolveAddressTest( //nolint:thelper
	tb testing.TB,
	tc resolveAddressTestCase,
	clientsFn func() (k8s.Clients, func(tb testing.TB)),
) {
	uri := &apis.URL{}
	clients, cleanup := clientsFn()
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
		mtBroker(namespace),
		channel(namespace),
	}
	for _, tc := range tcs {
		casefn(tc)
	}
}

func k8sService(namespace string) resolveAddressTestCase {
	svc := corev1.Service{
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
	}
	svc.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))
	return resolveAddressTestCase{
		name: "k8s service",
		matches: func(u *url.URL) bool {
			return u.String() ==
				fmt.Sprintf("http://k8s-hello.%s.svc.cluster.local/", namespace)
		},
		err:     nil,
		ref:     toTrackerRef(&svc),
		objects: []runtime.Object{&svc},
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
	m.url = clusterLocalURL
	publicURL := apis.HTTP(fmt.Sprintf(
		"%s-%s.apps.cloud.example.org", m.name, namespace))
	serviceURL := publicURL
	if clusterLocal {
		serviceURL = clusterLocalURL
	}
	ksvc := servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.name, Namespace: namespace, Labels: labels,
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
	}
	ksvc.SetGroupVersionKind(servingv1.SchemeGroupVersion.WithKind("Service"))
	return resolveAddressTestCase{
		name:    m.name,
		matches: m.matches,
		err:     nil,
		ref:     toTrackerRef(&ksvc),
		objects: []runtime.Object{&ksvc},
	}
}

func mtBroker(namespace string) resolveAddressTestCase {
	name := "test"
	u := apis.HTTP("broker-ingress.knative-eventing.svc.cluster.local")
	u.Path = fmt.Sprintf("/%s/%s", namespace, name)
	br := eventingv1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name: name, Namespace: namespace,
		},
		Status: eventingv1.BrokerStatus{
			Address: duckv1.Addressable{URL: u},
		},
	}
	br.SetGroupVersionKind(eventingv1.SchemeGroupVersion.WithKind("Broker"))
	return resolveAddressTestCase{
		name:    "mt-broker",
		err:     nil,
		ref:     toTrackerRef(&br),
		matches: equals(u),
		objects: []runtime.Object{&br},
	}
}

func channel(namespace string) resolveAddressTestCase {
	name := "test"
	u := apis.HTTP(
		fmt.Sprintf("%s-kn-channel.%s.svc.cluster.local", name, namespace))
	ch := messagingv1.Channel{
		ObjectMeta: metav1.ObjectMeta{
			Name: name, Namespace: namespace,
		},
		Status: messagingv1.ChannelStatus{
			ChannelableStatus: eventingduckv1.ChannelableStatus{
				AddressStatus: duckv1.AddressStatus{
					Address: &duckv1.Addressable{URL: u},
				},
			},
		},
	}
	ch.SetGroupVersionKind(messagingv1.SchemeGroupVersion.WithKind("Channel"))
	return resolveAddressTestCase{
		name:    "channel",
		err:     nil,
		ref:     toTrackerRef(&ch),
		matches: equals(u),
		objects: []runtime.Object{&ch},
	}
}

func toTrackerRef(accessor kmeta.Accessor) *tracker.Reference {
	gvk := accessor.GroupVersionKind()
	return &tracker.Reference{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Namespace:  accessor.GetNamespace(),
		Name:       accessor.GetName(),
		Selector:   nil,
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

func equals(u *apis.URL) func(url *url.URL) bool {
	return func(url *url.URL) bool {
		return u.String() == url.String()
	}
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

func fakeClients(tb testing.TB, tc resolveAddressTestCase) k8s.Clients {
	tb.Helper()
	return &tests.FakeClients{
		Objects: tc.objects,
		TB:      tb,
	}
}

type resolveAddressTestCase struct {
	name    string
	matches func(url *url.URL) bool
	err     error
	ref     *tracker.Reference
	objects []runtime.Object
}
