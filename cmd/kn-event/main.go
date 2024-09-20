package main

import (
	"github.com/wavesoftware/go-commandline"
	"knative.dev/kn-plugin-event/internal/cli"
)

func main() {
	commandline.New(new(cli.App)).ExecuteOrDie(cli.Options...)
}
