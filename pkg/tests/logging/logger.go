package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"knative.dev/pkg/logging"
)

func contextLogger(ctx context.Context) *zap.SugaredLogger {
	fallback := logging.FromContext(ctx)
	zlog, err := zap.NewDevelopment(loggerOptions()...)
	if err != nil {
		fallback.Fatal(err)
	}
	return zlog.Sugar()
}

func testLogger(t zaptest.TestingT) *zap.SugaredLogger {
	return zaptest.NewLogger(t,
		zaptest.WrapOptions(loggerOptions()...)).Sugar()
}

func loggerOptions() []zap.Option {
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.WarnLevel)),
	}
}
