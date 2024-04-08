package test

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/tracker"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

const (
	// HTTPPort is 80.
	HTTPPort = 80

	// AlternativeHTTPPort is 8080.
	AlternativeHTTPPort = 8080
)

var (
	ErrNotEqual    = errors.New("not equal")
	ErrDontContain = errors.New("don't contain")
)

func ResolveAddressTestCases(namespace string, casefn func(tc ResolveAddressTestCase)) {
	tcs := []ResolveAddressTestCase{
		k8sService(namespace),
		knService(namespace),
		mtBroker(namespace),
		channel(namespace),
	}
	for _, tc := range tcs {
		casefn(tc)
	}
}

// EnsureResolveAddress thelper lint skipped for greater visibility of
// failure location.
func EnsureResolveAddress( //nolint:thelper
	tb testing.TB,
	tc ResolveAddressTestCase,
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
	assert.NilError(tb, tc.matches(u))
}

type ResolveAddressTestCase struct {
	Name    string
	matches func(url *url.URL) error
	err     error
	ref     *tracker.Reference
	Objects []runtime.Object
}

func k8sService(namespace string) ResolveAddressTestCase {
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-hello",
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: HTTPPort,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: AlternativeHTTPPort,
				},
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	svc.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))
	want := apis.HTTP(fmt.Sprintf("k8s-hello.%s.svc.cluster.local", namespace))
	return ResolveAddressTestCase{
		Name:    "k8s-service",
		matches: equals(want),
		err:     nil,
		ref:     toTrackerRef(&svc),
		Objects: []runtime.Object{&svc},
	}
}

func knService(namespace string) ResolveAddressTestCase {
	m := matcher{
		name:      "kn-hello",
		namespace: namespace,
	}
	labels := map[string]string{}
	clusterLocalURL := apis.HTTP(fmt.Sprintf(
		"%s.%s.svc.cluster.local", m.name, namespace))
	m.url = clusterLocalURL
	publicURL := apis.HTTP(fmt.Sprintf(
		"%s-%s.apps.cloud.example.org", m.name, namespace))
	serviceURL := publicURL
	ksvc := servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.name, Namespace: namespace, Labels: labels,
		},
		Spec: ksvcSpec("quay.io/openshift-knative/helloworld-go:multiarch"),
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
	return ResolveAddressTestCase{
		Name:    m.name,
		matches: m.matches,
		err:     nil,
		ref:     toTrackerRef(&ksvc),
		Objects: []runtime.Object{&ksvc},
	}
}

func mtBroker(namespace string) ResolveAddressTestCase {
	name := "test"
	u := apis.HTTP("broker-ingress.knative-eventing.svc.cluster.local")
	u.Path = fmt.Sprintf("/%s/%s", namespace, name)
	br := eventingv1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name: name, Namespace: namespace,
		},
		Status: eventingv1.BrokerStatus{
			AddressStatus: duckv1.AddressStatus{
				Address: &duckv1.Addressable{URL: u},
			},
		},
	}
	br.SetGroupVersionKind(eventingv1.SchemeGroupVersion.WithKind("Broker"))
	return ResolveAddressTestCase{
		Name:    "mt-broker",
		err:     nil,
		ref:     toTrackerRef(&br),
		matches: equals(u),
		Objects: []runtime.Object{&br},
	}
}

func channel(namespace string) ResolveAddressTestCase {
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
	return ResolveAddressTestCase{
		Name:    "channel",
		err:     nil,
		ref:     toTrackerRef(&ch),
		matches: equals(u),
		Objects: []runtime.Object{&ch},
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
	name      string
	namespace string
}

func (m matcher) matches(u *url.URL) error {
	if !strings.Contains(u.Host, m.name) {
		return fmt.Errorf("%w: expect %v to contain %v", ErrDontContain, u, m.name)
	}
	if !strings.Contains(u.Host, m.namespace) {
		return fmt.Errorf("%w: expect %v to contain %v", ErrDontContain, u, m.namespace)
	}
	return check(u,
		m.containsName,
		m.containsNamespace,
	)
}

func (m matcher) containsName(u *url.URL) error {
	return hostContains(u, m.name)
}

func (m matcher) containsNamespace(u *url.URL) error {
	return hostContains(u, m.namespace)
}

func hostContains(u *url.URL, needle string) error {
	if !strings.Contains(u.Host, needle) {
		return fmt.Errorf("%w: expect %v to contain %#v",
			ErrDontContain, u, needle)
	}
	return nil
}

func check(u *url.URL, fns ...func(*url.URL) error) error {
	for _, fn := range fns {
		err := fn(u)
		if err != nil {
			return err
		}
	}
	return nil
}

func equals(u *apis.URL) func(url *url.URL) error {
	return func(url *url.URL) error {
		if u.String() != url.String() {
			return fmt.Errorf("%w: got %v, want %v", ErrNotEqual, url, u)
		}
		return nil
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
