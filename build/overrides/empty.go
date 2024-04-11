package overrides

import "knative.dev/toolbox/magetasks/config"

// List holds overrides by which downstream forks could influence the build
// process.
//
// goland:noinspection GoUnusedGlobalVariable
var List []config.Configurator //nolint:gochecknoglobals
