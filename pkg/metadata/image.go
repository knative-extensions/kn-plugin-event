package metadata

import (
	pkgimage "github.com/wavesoftware/go-magetasks/pkg/image"
)

var (
	// Image holds information about companion image reference.
	Image = "" //nolint:gochecknoglobals
	// ImageBasename holds a basename of a image, so the development reference
	// could be built from it.
	ImageBasename = "" //nolint:gochecknoglobals
	// ImageBasenameSeparator holds a separator between image basename and name.
	ImageBasenameSeparator = "/" //nolint:gochecknoglobals
)

// ResolveImage will try to resolve the image reference from set values. If
// Image is given it will be used, otherwise the ImageBasename and Version will
// be used.
func ResolveImage() string {
	//goland:noinspection GoBoolExpressions
	if Image == "" {
		return pkgimage.FloatToRelease(
			ImageBasename, "kn-event-sender", ImageBasenameSeparator, Version,
			pkgimage.FloatDirectionDown)
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

// ImageBasenameSeparatorPath return a path to the image basename separator
// variable.
func ImageBasenameSeparatorPath() string {
	return importPath("ImageBasenameSeparator")
}
