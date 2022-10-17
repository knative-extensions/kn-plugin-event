package version

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
)

var (
	// ErrIsNotValid when version data is malformed.
	ErrIsNotValid = errors.New("is not valid")

	// ErrVersionIsNotValid when given version is not semantic.
	ErrVersionIsNotValid = fmt.Errorf("version %w", ErrIsNotValid)

	// ErrRangeIsNotValid when given range is not valid.
	ErrRangeIsNotValid = fmt.Errorf("range %w", ErrIsNotValid)

	// ErrVersionOutsideOfRange when given version is outside of range.
	ErrVersionOutsideOfRange = fmt.Errorf("version %w for given range", ErrIsNotValid)
)

// Resolver will resolve version string, and tell is that the latest artifact.
type Resolver interface {
	// Version returns the version string.
	Version() string
	// IsLatest tells if the version is the latest one within given version range.
	IsLatest(versionRange string) (bool, error)
}

// CompatibleRanges will resolve compatible ranges from a version resolver.
func CompatibleRanges(resolver Resolver) ([]string, error) {
	v := resolver.Version()
	prefix := ""
	if strings.HasPrefix(v, "v") {
		prefix = "v"
	}
	sv, err := semver.ParseTolerant(v)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVersionIsNotValid, err)
	}
	if len(sv.Pre) == 0 && len(sv.Build) == 0 {
		return collect(
			maybeMajorRange(prefix, sv, resolver),
			maybeMinorRange(prefix, sv, resolver),
		)
	}
	return []string{}, nil
}

type either struct {
	slice []string
	err   error
}

func collect(eithers ...either) ([]string, error) {
	ret := make([]string, 0, len(eithers))
	for _, e := range eithers {
		if e.err != nil {
			return nil, e.err
		}
		ret = append(ret, e.slice...)
	}
	return ret, nil
}

type ranges struct {
	constraint, short string
}

func maybeMajorRange(prefix string, version semver.Version, resolver Resolver) either {
	current := fmt.Sprintf("%d.0.0", version.Major)
	next := fmt.Sprintf("%d.0.0", version.Major+1)
	return maybeRange(ranges{
		constraint: fmt.Sprintf(">= %s < %s", current, next),
		short:      fmt.Sprintf("%s%d", prefix, version.Major),
	}, resolver)
}

func maybeMinorRange(prefix string, version semver.Version, resolver Resolver) either {
	current := fmt.Sprintf("%d.%d.0", version.Major, version.Minor)
	next := fmt.Sprintf("%d.%d.0", version.Major, version.Minor+1)
	return maybeRange(ranges{
		constraint: fmt.Sprintf(">= %s < %s", current, next),
		short:      fmt.Sprintf("%s%d.%d", prefix, version.Major, version.Minor),
	}, resolver)
}

func maybeRange(r ranges, resolver Resolver) either {
	ok, err := resolver.IsLatest(r.constraint)
	if err != nil {
		return either{err: err}
	}
	if ok {
		return either{slice: []string{r.short}}
	}
	return either{}
}
