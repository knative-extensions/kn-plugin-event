package ics

import (
	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/pkg/binding"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/pkg/logging"
)

// Options to override the commandline for testing purposes.
var Options []commandline.Option //nolint:gochecknoglobals

type App struct {
	k8s.Params
}

func (a App) Command() *cobra.Command {
	c := &cobra.Command{
		Use:           "ics",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          a.run,
	}
	c.SetContext(cli.InitialContext())
	c.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		cli.SetupContext(cmd, zapcore.DebugLevel)
	}
	c.PersistentPostRunE = func(cmd *cobra.Command, _ []string) error {
		closer := outlogging.LogFileCloserFrom(cmd.Context())
		// ensure to close the log file
		return closer()
	}
	a.SetGlobalFlags(c.PersistentFlags())
	return c
}

func (a App) run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log := logging.FromContext(ctx)
	err := binding.IcsApp().SendFromEnv(ctx, a.Parse())
	if err != nil {
		log.Error(zap.Error(err))
	}
	return err
}

var _ commandline.CobraProvider = new(App)
