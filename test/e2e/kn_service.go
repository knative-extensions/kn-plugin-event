//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/kmeta"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/k8s"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingv1clientset "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

const watholaForwarderPackage = "knative.dev/eventing/test/test_images/wathola-forwarder"

// RegisterPackages will register packages to be built into test images.
func RegisterPackages() {
	environment.RegisterPackage(watholaForwarderPackage)
}

// SendEventToKnService returns a feature.Feature that verifies the kn-event
// can send to Knative service.
func SendEventToKnService() *feature.Feature {
	return sendEventFeature("ToKnativeService", sendEventOptions{
		sink: func(sinkName string) string {
			return watholaForwarder{sinkName}.sinkSpec()
		},
		steps: []step{
			func(f *feature.Feature, sinkName string) *feature.Feature {
				f.Setup("deploy wathola-forwarder",
					watholaForwarder{sinkName}.step)
				return f
			},
		},
	})
}

type watholaForwarder struct {
	sinkName string
}

func (wf watholaForwarder) step(ctx context.Context, t feature.T) {
	wf.deployConfigMap(ctx, t)
	wf.deployKnService(ctx, t)
}

//go:embed wathola-forwarder.toml
var watholaForwarderConfigTemplate embed.FS

func (wf watholaForwarder) deployConfigMap(ctx context.Context, t feature.T) {
	env := environment.FromContext(ctx)
	ns := env.Namespace()
	tmpl, err := template.ParseFS(watholaForwarderConfigTemplate,
		"wathola-forwarder.toml")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, struct {
		Sink *apis.URL
	}{
		apis.HTTP(fmt.Sprintf("%s.%s.svc", wf.sinkName, ns)),
	}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	kube := kubeclient.Get(ctx)
	configMaps := kube.CoreV1().ConfigMaps(ns)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: wf.name(), Namespace: ns},
		Data: map[string]string{
			path.Base(wf.configPath()): buff.String(),
		},
	}
	created, err := configMaps.Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	env.Reference(kmeta.ObjectReference(created))
}

func (wf watholaForwarder) deployKnService(ctx context.Context, t feature.T) {
	env := environment.FromContext(ctx)
	ns := env.Namespace()
	rest := injection.GetConfig(ctx)
	ksvcs := servingv1clientset.NewForConfigOrDie(rest).Services(ns)
	image := fmt.Sprintf("ko://%s", watholaForwarderPackage)
	if replaced, found := env.Images()[image]; found {
		image = replaced
	}
	const readyPath = "/healthz"
	ksvc := &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: wf.name(), Namespace: ns},
		Spec: servingv1.ServiceSpec{
			ConfigurationSpec: servingv1.ConfigurationSpec{
				Template: servingv1.RevisionTemplateSpec{
					Spec: servingv1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Image: image,
								ReadinessProbe: &corev1.Probe{
									Handler: corev1.Handler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: readyPath,
										},
									},
								},
								VolumeMounts: []corev1.VolumeMount{{
									Name:      wf.name(),
									MountPath: path.Dir(wf.configPath()),
									ReadOnly:  true,
								}},
							}},
							Volumes: []corev1.Volume{{
								Name: wf.name(),
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: wf.name(),
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	}
	created, err := ksvcs.Create(ctx, ksvc, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	ref := kmeta.ObjectReference(created)
	env.Reference(ref)
	if err = k8s.WaitForServiceReady(ctx, t, wf.name(), readyPath); err != nil {
		t.Fatal(errors.WithStack(err))
	}
}

func (wf watholaForwarder) name() string {
	return wf.sinkName + "-wathola-forwarder"
}

func (wf watholaForwarder) sinkSpec() string {
	return fmt.Sprintf("Service:serving.knative.dev/v1:%s", wf.name())
}

func (wf watholaForwarder) configPath() string {
	homedir := "/home/nonroot"
	if homedirEnv, ok := os.LookupEnv("KN_PLUGIN_EVENT_WATHOLA_HOMEDIR"); ok {
		homedir = homedirEnv
	}
	return fmt.Sprintf("%s/.config/wathola/config.toml", homedir)
}
