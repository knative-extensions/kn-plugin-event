package reference_test

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/logging"
	"knative.dev/reconciler-test/pkg/feature"
)

func testContext(t zaptest.TestingT) context.Context {
	log := zaptest.NewLogger(t, zaptest.WrapOptions(
		zap.AddCaller(),
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.WarnLevel)),
	)).Sugar()
	return logging.WithLogger(context.TODO(), log)
}

func meta(name, ns string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      feature.MakeRandomK8sName(name),
		Namespace: feature.MakeRandomK8sName(ns),
	}
}
