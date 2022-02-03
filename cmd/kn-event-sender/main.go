package main

import (
	"context"
	"os"

	"go.uber.org/zap"
	"knative.dev/kn-plugin-event/pkg/cli/retcode"
	"knative.dev/kn-plugin-event/pkg/configuration"
	"knative.dev/pkg/logging"
)

// ExitFunc will be used to exit Go process in case of error.
var ExitFunc = os.Exit // nolint:gochecknoglobals

func main() {
	logger := createLogger()
	app := configuration.CreateIcs()
	if err := app.SendFromEnv(); err != nil {
		logger.Error(err)
		ExitFunc(retcode.Calc(err))
	}
}

// TestMain is used by tests.
//goland:noinspection GoUnusedExportedFunction
func TestMain() { //nolint:deadcode
	main()
}

func createLogger() *zap.SugaredLogger {
	return logging.
		FromContext(context.TODO()).
		With("env", os.Environ())
}
