//go:build e2e
// +build e2e

package e2e

import (
	"fmt"

	"knative.dev/reconciler-test/pkg/feature"
)

// SystemUnderTest is a part of cluster we are testing the event propagation
// though.
type SystemUnderTest interface {
	Name() string
	Deploy(feature *feature.Feature, sinkName string) Sink
}

// Sink represents a parameter in format acceptable by the `--to` option.
type Sink interface {
	fmt.Stringer
}

type sinkFn func() string

func (s sinkFn) String() string {
	return s()
}
