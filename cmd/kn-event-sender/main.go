package main

import (
	"fmt"
	"os"

	"knative.dev/kn-plugin-event/pkg/cli/retcode"
	"knative.dev/kn-plugin-event/pkg/configuration"
)

// ExitFunc will be used to exit Go process in case of error.
var ExitFunc = os.Exit // nolint:gochecknoglobals

func main() {
	app := configuration.CreateIcs()
	if err := app.SendFromEnv(); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		ExitFunc(retcode.Calc(err))
	}
}

// TestMain is used by tests.
//goland:noinspection GoUnusedExportedFunction
func TestMain() { //nolint:deadcode
	main()
}
