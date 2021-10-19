package ldflags

import (
	"fmt"
	"strings"

	"github.com/wavesoftware/go-magetasks/config"
)

// Builder builds the LD flags by adding values resolvers.
type Builder interface {
	// Add a name and a resolver to the builder.
	Add(name string, resolver config.Resolver) Builder
	// Build into slice args for ldflags.
	Build() []string
	// BuildOnto provided args.
	BuildOnto(args []string) []string
}

// NewBuilder creates a new builder.
func NewBuilder() Builder {
	return &defaultBuilder{
		resolvers: make(map[string]config.Resolver),
	}
}

type defaultBuilder struct {
	resolvers map[string]config.Resolver
}

func (d *defaultBuilder) Add(name string, resolver config.Resolver) Builder {
	d.resolvers[name] = resolver
	return d
}

func (d *defaultBuilder) Build() []string {
	collected := make([]string, 0, len(d.resolvers))
	if len(d.resolvers) == 0 {
		return collected
	}
	for name, resolver := range d.resolvers {
		collected = append(collected, fmt.Sprintf("-X %s=%s", name, resolver()))
	}
	return collected
}

func (d *defaultBuilder) BuildOnto(args []string) []string {
	return append(args, "-ldflags", strings.Join(d.Build(), " "))
}
