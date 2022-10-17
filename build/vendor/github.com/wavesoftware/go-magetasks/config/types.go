package config

import (
	"context"

	"github.com/fatih/color"
	"github.com/wavesoftware/go-magetasks/pkg/version"
)

// Config holds configuration information.
type Config struct {
	// ProjectDir holds a path to the project.
	ProjectDir string

	// BuildDirPath holds a build directory path.
	BuildDirPath []string

	// Version contains the version information.
	*Version

	// MageTag holds default mage tag settings.
	MageTag

	// Dependencies will hold additional Golang dependencies that needs to be
	// installed before running tasks.
	Dependencies Dependencies

	// Artifacts holds a list of artifacts to be built.
	Artifacts []Artifact

	// Builders holds a list of Builder's to be used for building project
	// artifacts. If none is configured, default ones will be used.
	Builders []Builder

	// Publishers holds a list of Publisher's to be used for publishing project
	// artifacts. If none is configured, default ones will be used.
	Publishers []Publisher

	// Cleaning additional cleaning tasks.
	Cleaning []Task

	// Checks holds a list of checks to perform.
	Checks []Task

	// BuildVariables holds extra list of variables to be passed to Golang's
	// ldflags.
	BuildVariables

	// Overrides holds a list of overrides of this configuration.
	Overrides []Configurator

	// context.Context is standard Golang context.
	context.Context //nolint:containedctx
}

// Notifier can notify of a pending status of long task.
type Notifier interface {
	Notify(status string)
}

// Builder builds an artifact.
type Builder interface {
	Accepts(artifact Artifact) bool
	Build(artifact Artifact, notifier Notifier) Result
}

// Publisher publishes artifacts to a remote site.
type Publisher interface {
	Accepts(artifact Artifact) bool
	Publish(artifact Artifact, notifier Notifier) Result
}

// Resolver is a func that resolves to a string.
type Resolver func() string

// StaticResolver can be used to create a Resolver from static data.
func StaticResolver(value string) Resolver {
	return func() string {
		return value
	}
}

// MageTag holds a mage tag.
type MageTag struct {
	Color color.Attribute
	Label string
}

// Result hold a result of an Artifact build or publish.
type Result struct {
	Error error
	Info  map[string]interface{}
}

// Failed returns true if the artifact processing failed.
func (r Result) Failed() bool {
	return r.Error != nil
}

// Artifact represents a thing that can be built and published.
type Artifact interface {
	GetType() string
	GetName() string
}

// ResultKey is used to store results of a build or publish.
type ResultKey string

// Version specifies the version information and how to set it into variable at
// compile time.
type Version struct {
	Path string
	version.Resolver
}

// Metadata holds additional contextual information.
type Metadata struct {
	Name string
	BuildVariables
}

func (m Metadata) GetName() string {
	return m.Name
}

// Task is a custom function that will be used in the build.
type Task struct {
	Name      string
	Operation func(Notifier) error
	Overrides []Configurator
}

// Configurator will configure project.
type Configurator interface {
	Configure(cfg Configurable)
}

// Configurable will allow changes of the Config structure.
type Configurable interface {
	Config() *Config
}

// Configured represents a configured project.
type Configured interface{}
