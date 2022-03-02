package images

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

var (
	// ErrNotFound is returned when there is no environment variable set for
	// given KO path.
	ErrNotFound     = errors.New("expected environment variable not found")
	nonAlphaNumeric = regexp.MustCompile("[^A-Z0-9]+")
)

// EnvironmentalBasedResolver will try to resolve the images from prefixed
// environment variables.
type EnvironmentalBasedResolver struct {
	Prefix string
}

func (c *EnvironmentalBasedResolver) Applicable() bool {
	prefix := c.normalizedPrefix()
	for _, environment := range os.Environ() {
		const equalitySignParts = 2
		parts := strings.SplitN(environment, "=", equalitySignParts)
		key := parts[0]
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

func (c *EnvironmentalBasedResolver) Resolve(kopath string) (name.Reference, error) {
	prefix := c.normalizedPrefix()
	shortName := normalize(path.Base(kopath))
	key := fmt.Sprintf("%s_%s", prefix, shortName)
	val, ok := os.LookupEnv(key)
	if !ok {
		return nil, fmt.Errorf("%w: '%s' - kopath: %s", ErrNotFound, key, kopath)
	}
	ref, err := name.ParseReference(val)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return ref, nil
}

func (c *EnvironmentalBasedResolver) normalizedPrefix() string {
	return normalize(c.Prefix)
}

func normalize(in string) string {
	return strings.Trim(
		nonAlphaNumeric.ReplaceAllString(strings.ToUpper(in), "_"),
		"_",
	)
}

var _ Resolver = &EnvironmentalBasedResolver{}
