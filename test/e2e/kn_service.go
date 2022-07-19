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
	"knative.dev/kn-plugin-event/pkg/tests/reference"
	"knative.dev/pkg/apis"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/k8s"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingv1clientset "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// SendEventToKnService returns a feature.Feature that verifies the kn-event
// can send to Knative service.
func SendEventToKnService() *feature.Feature {
	return SendEventFeature(knServiceSut{})
}

type knServiceSut struct{}

func (k knServiceSut) Name() string {
	return "KnativeService"
}

func (k knServiceSut) Deploy(f *feature.Feature, sinkName string) Sink {
	wf := watholaForwarder{sinkName}
	f.Setup("Deploy KnService", wf.step)
	return wf.sink()
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
	content := buff.String()
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: wf.name(), Namespace: ns},
		Data: map[string]string{
			path.Base(wf.configPath()): content,
		},
	}
	log := logging.FromContext(ctx).
		With(json("meta", cm.ObjectMeta))
	log.Info("Deploying ConfigMap")
	if _, err = configMaps.Create(ctx, cm, metav1.CreateOptions{}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	env.Reference(reference.FromConfigMap(ctx, cm))
	log.With(json("content", content)).Info("ConfigMap is deployed")
}

func (wf watholaForwarder) deployKnService(ctx context.Context, t feature.T) {
	env := environment.FromContext(ctx)
	ns := env.Namespace()
	rest := injection.GetConfig(ctx)
	ksvcs := servingv1clientset.NewForConfigOrDie(rest).Services(ns)
	image := wf.image(ctx, t)
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
									ProbeHandler: corev1.ProbeHandler{
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
	log := logging.FromContext(ctx).
		With(json("meta", ksvc.ObjectMeta))
	log.With(json("spec", ksvc.Spec.Template.Spec.PodSpec)).
		Info("Deploying KnService")
	if _, err := ksvcs.Create(ctx, ksvc, metav1.CreateOptions{}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	ref := reference.FromKnativeService(ctx, ksvc)
	env.Reference(ref)
	k8s.WaitForReadyOrDoneOrFail(ctx, t, ref)
	log.Info("KnService is ready")
}

func (wf watholaForwarder) image(ctx context.Context, t feature.T) string {
	image := WatholaForwarderImageFromContext(ctx)
	if images, err := environment.ProduceImages(ctx); err != nil {
		t.Fatal(errors.WithStack(err))
	} else {
		if replaced, found := images[image]; found {
			image = replaced
		}
	}
	return image
}

func (wf watholaForwarder) name() string {
	return wf.sinkName + "-ksvc"
}

func (wf watholaForwarder) sink() Sink {
	return sinkFn(func() string {
		return fmt.Sprintf("Service:%s:%s",
			servingv1.SchemeGroupVersion, wf.name())
	})
}

func (wf watholaForwarder) configPath() string {
	homedir := "/home/nonroot"
	if homedirEnv, ok := os.LookupEnv("KN_PLUGIN_EVENT_WATHOLA_HOMEDIR"); ok {
		homedir = homedirEnv
	}
	return fmt.Sprintf("%s/.config/wathola/config.toml", homedir)
}
