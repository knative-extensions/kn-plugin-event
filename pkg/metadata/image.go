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

// ImagePath return a path to the image variable.
func ImagePath() string {
	return importPath("Image")
}

// ImageBasenamePath return a path to the image basename variable.
func ImageBasenamePath() string {
	return importPath("ImageBasename")
}
