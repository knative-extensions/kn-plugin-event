package main

import (
	"log"
	"os"

	"go.uber.org/zap"
	"knative.dev/kn-plugin-event/pkg/cli/retcode"
	"knative.dev/kn-plugin-event/pkg/configuration"
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
	zapl, err := zap.NewProduction()
	if err != nil {
		log.Println(err)
		ExitFunc(retcode.Calc(err))
		return nil
	}
	return zapl.Sugar().With("env", os.Environ())
}
