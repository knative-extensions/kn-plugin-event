package buildvars

import "github.com/wavesoftware/go-magetasks/config"

// Operation performs some operation on Builder and return modified one.
type Operation func(Builder) Builder

type Operator interface {
	Operation() Operation
}

// Assemble will assemble a set of operations into the build variables.
func Assemble(operators []Operator) config.BuildVariables {
	b := Builder{}
	for _, operator := range operators {
		b = operator.Operation()(b)
	}
	return b.Build()
}
