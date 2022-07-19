package reference_test

import (
	"context"

	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/logging"
)

func testContext(t zaptest.TestingT) context.Context {
	return logging.WithTestLogger(logging.NewContext(context.TODO()), t)
}

func meta(name, ns string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      feature.MakeRandomK8sName(name),
		Namespace: feature.MakeRandomK8sName(ns),
	}
}
