package ics

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"knative.dev/kn-plugin-event/pkg/configuration"
	"knative.dev/kn-plugin-event/pkg/system"
	"knative.dev/pkg/logging"
)

// Options to override the commandline for testing purposes.
var Options []commandline.Option //nolint:gochecknoglobals

type App struct{}

func (a App) Command() *cobra.Command {
	return &cobra.Command{
		Use:           "ics",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          a.run,
	}
}

func (a App) run(cmd *cobra.Command, _ []string) error {
	env := withLogger(cmd)
	log := logging.FromContext(env.Context())
	defer func(log *zap.SugaredLogger) {
		_ = log.Sync()
	}(log)
	err := configuration.CreateIcs(env).SendFromEnv()
	if err != nil {
		log.Error(zap.Error(err))
	}
	return err //nolint:wrapcheck
}

var _ commandline.CobraProvider = new(App)

func withLogger(env system.Environment) system.Environment {
	ctx := env.Context()
	ctx = logging.WithLogger(ctx, createLogger(env))
	return system.WithContext(ctx, env)
}

func createLogger(env system.Environment) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	sink := zapcore.AddSync(env.OutOrStdout())
	zcore := zapcore.NewCore(encoder, sink, cfg.Level)
	return zap.New(zcore).
		With(zap.Strings("env", os.Environ())).
		Sugar()
}
