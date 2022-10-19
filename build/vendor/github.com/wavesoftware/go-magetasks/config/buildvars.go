package config

// BuildVariables will be passed to a Golang's ldflags for variable injection.
type BuildVariables map[string]Resolver
