package image

import (
	"errors"

	"github.com/wavesoftware/go-magetasks/pkg/version"
)

// Tags will return a list of tags for an OCI image based on the version
// information given by resolver.
func Tags(resolver version.Resolver) ([]string, error) {
	ranges, err := version.CompatibleRanges(resolver)
	if err != nil {
		if !errors.Is(err, version.ErrVersionIsNotValid) {
			return nil, err
		}
		ranges = make([]string, 0)
	}
	tags := append([]string{resolver.Version()}, ranges...)
	var latest bool
	latest, err = resolver.IsLatest(version.AnyVersion)
	if err != nil && !errors.Is(err, version.ErrVersionIsNotValid) {
		return nil, err
	}
	if latest {
		tags = append(tags, "latest")
	}
	return tags, nil
}
