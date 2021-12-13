package artifact

import "github.com/wavesoftware/go-magetasks/pkg/artifact/platform"

// Platform to built binary for.
type Platform struct {
	platform.OS
	platform.Architecture
}
