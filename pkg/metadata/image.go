package metadata

import "fmt"

var (
	// Image holds information about companion image reference.
	Image = "" //nolint:gochecknoglobals
	// ImageBasename holds a basename of a image, so the development reference
	// could be built from it.
	ImageBasename = "" //nolint:gochecknoglobals
)

// ResolveImage will try to resolve the image reference from set values. If
// Image is given it will be used, otherwise the ImageBasename and Version will
// be used.
func ResolveImage() string {
	//goland:noinspection GoBoolExpressions
	if Image == "" {
		return fmt.Sprintf("%s/kn-event-sender:%s", ImageBasename, Version)
	}
	return Image
}
