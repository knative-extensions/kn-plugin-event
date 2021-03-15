package main

import (
	"fmt"
	"os"

	"github.com/cardil/kn-event/internal/cli/retcode"
	"github.com/cardil/kn-event/internal/configuration"
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
