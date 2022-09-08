package main

import (
	"github.com/wavesoftware/go-commandline"
	"knative.dev/kn-plugin-event/internal/ics"
)

func main() {
	commandline.New(ics.App{}).ExecuteOrDie(ics.Options...)
}

// TestMain is used by tests.
func TestMain() { //nolint:deadcode
	main()
}
