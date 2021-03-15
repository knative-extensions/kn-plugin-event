package k8s_test

import (
	"errors"
	"testing"

	"github.com/cardil/kn-event/internal/k8s"
	"github.com/cardil/kn-event/internal/tests"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/tracker"
)

func TestResolveAddress(t *testing.T) {
	testCases := []resolveAddressTestCase{
		passingService(t),
	}
	for i := range testCases {
		tt := testCases[i]
		t.Run(tt.name, func(t *testing.T) {
			// FIXME: can't discover GVK using fake clients - endless loop
			t.Skip("can't discover GVK using fake clients - endless loop")
			resolver := k8s.CreateAddressResolver(tt.args.clients)
			uri, err := apis.ParseURL("/")
			assert.NoError(t, err)
			u, err := resolver.ResolveAddress(tt.args.ref, uri)
			if errors.Is(err, tt.err) {
				t.Errorf("CreateAddressResolver() err = %v, tt.err = %v", err, tt.err)
				return
			}
			assert.Equal(t, tt.wantURL, u.String())
		})
	}
}

func passingService(t *testing.T) resolveAddressTestCase {
	t.Helper()
	clients := &tests.FakeClients{Objects: []runtime.Object{
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hello",
				Namespace: "demo",
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
		},
	}}
	return resolveAddressTestCase{
		name: "passing",
		args: args{
			clients: clients,
			ref: &tracker.Reference{
				APIVersion: "v1",
				Kind:       "Service",
				Namespace:  "demo",
				Name:       "hello",
				Selector:   nil,
			},
		},
		wantURL: "hello.demo.svc",
		err:     nil,
	}
}

type resolveAddressTestCase struct {
	name    string
	args    args
	wantURL string
	err     error
}

type args struct {
	clients k8s.Clients
	ref     *tracker.Reference
}
