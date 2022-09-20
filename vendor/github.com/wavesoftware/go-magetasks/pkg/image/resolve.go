package image

import (
	"fmt"

	"github.com/blang/semver/v4"
)

// VersionPrefix used to prefix versions.
var VersionPrefix = "v"

// FloatDirection defines the way to float a non-release version.
type FloatDirection int

const (
	// FloatDirectionUp means the minor version will be incremented to find
	// compatible range ,effectively meaning a next minor release.
	FloatDirectionUp FloatDirection = iota
	// FloatDirectionDown means the minor version will be left intact, but patch
	// number will be removed, effectively meaning latest version from current
	// minor release.
	FloatDirectionDown
)

// FloatToRelease will build a full image name from basename, name, separator,
// version parts given as arguments. If version is a non-release it will be
// floated either up or down depending on direction argument. Floating up means
// to increase the minor number by 1. Floating down means leaving minor number
// as it was.
func FloatToRelease(basename, name, separator, version string, direction FloatDirection) string {
	sver, err := semver.ParseTolerant(version)
	if err == nil {
		version = fmt.Sprintf("%s%d.%d.%d", VersionPrefix, sver.Major, sver.Minor, sver.Patch)
		if len(sver.Pre) > 0 || len(sver.Build) > 0 {
			// non release image
			major := sver.Major
			minor := sver.Minor
			if direction == FloatDirectionUp {
				minor++
			}
			version = fmt.Sprintf("%s%d.%d", VersionPrefix, major, minor)
		}
	}
	return fmt.Sprintf("%s%s%s:%s", basename, separator, name, version)
}
