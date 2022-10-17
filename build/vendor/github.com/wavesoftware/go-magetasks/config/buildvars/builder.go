package buildvars

import "github.com/wavesoftware/go-magetasks/config"

// Builder for a build variables.
type Builder struct {
	bv config.BuildVariables
}

// Add a key/value pair.
func (b Builder) Add(key string, resolver config.Resolver) Builder {
	b.ensureBuildVariables()
	b.bv[key] = resolver
	return b
}

// ConditionallyAdd a key/value pair if cnd is true.
func (b Builder) ConditionallyAdd(cnd func() bool, key string, resolver config.Resolver) Builder {
	if cnd() {
		return b.Add(key, resolver)
	}
	return b
}

// Build a build variables instance.
func (b Builder) Build() config.BuildVariables {
	b.ensureBuildVariables()
	return b.bv
}

func (b *Builder) ensureBuildVariables() {
	if b.bv == nil {
		b.bv = make(config.BuildVariables)
	}
}
