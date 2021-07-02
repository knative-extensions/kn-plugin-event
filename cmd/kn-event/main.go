package main

import "knative.dev/kn-plugin-event/cmd/kn-event/cmd"

// Suppress global check for testing purposes.
var mainCmd = &cmd.Cmd{} //nolint:gochecknoglobals

func main() {
	mainCmd.Execute()
}
