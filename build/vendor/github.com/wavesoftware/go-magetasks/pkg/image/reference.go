package image

import (
	"os"
	"strings"

	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/config/buildvars"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
)

// InfluenceableReference defines an image reference that can be influenced by
// environment variables. To disable referencing images by they built hash,
// caller should set variable `DONT_REFERENCE_IMAGE_BY_DIGEST=true` - images
// will use version and basename to resolve full OCI image name.
//
// By setting variable defined in EnvVariable field, caller can override the
// full OCI name of the image.
type InfluenceableReference struct {
	Path        string
	EnvVariable string
	artifact.Image
}

func (r InfluenceableReference) Operation() buildvars.Operation {
	return func(builder buildvars.Builder) buildvars.Builder {
		return builder.ConditionallyAdd(
			referenceImageByDigest,
			r.Path,
			environmentOverrideImageReference(r.EnvVariable,
				artifact.ImageReferenceOf(r.Image)),
		)
	}
}

func environmentOverrideImageReference(
	envVariable string, fallbackResolver config.Resolver,
) config.Resolver {
	return func() string {
		if val, ok := os.LookupEnv(envVariable); ok {
			return val
		}
		return fallbackResolver()
	}
}

func dontReferenceImageByDigest() bool {
	if val, ok := os.LookupEnv("DONT_REFERENCE_IMAGE_BY_DIGEST"); ok {
		return strings.ToLower(val) == "true"
	}
	return false
}

func referenceImageByDigest() bool {
	return !dontReferenceImageByDigest()
}
