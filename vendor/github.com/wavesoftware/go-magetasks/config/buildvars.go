package config

// BuildVariables will be passed to a Golang's ldflags for variable injection.
type BuildVariables map[string]Resolver

// BuildVariablesBuilder for a build variables.
type BuildVariablesBuilder interface {
	// Add a key/value pair.
	Add(key string, resolver Resolver) BuildVariablesBuilder
	// ConditionallyAdd a key/value pair if cnd is true.
	ConditionallyAdd(cnd func() bool, key string, resolver Resolver) BuildVariablesBuilder
	// Build a build variables instance.
	Build() BuildVariables
}

// NewBuildVariablesBuilder creates a new BuildVariablesBuilder.
func NewBuildVariablesBuilder() BuildVariablesBuilder {
	return &defaultBuilder{
		bv: make(BuildVariables),
	}
}

type defaultBuilder struct {
	bv BuildVariables
}

func (d defaultBuilder) Add(key string, resolver Resolver) BuildVariablesBuilder {
	d.bv[key] = resolver
	return d
}

func (d defaultBuilder) ConditionallyAdd(cnd func() bool, key string, resolver Resolver) BuildVariablesBuilder {
	if cnd() {
		return d.Add(key, resolver)
	}
	return d
}

func (d defaultBuilder) Build() BuildVariables {
	return d.bv
}
