package main

import (
	"github.com/wavesoftware/go-commandline"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
)

func main() {
	commandline.New(new(cmd.App)).ExecuteOrDie(cmd.Options...)
}
