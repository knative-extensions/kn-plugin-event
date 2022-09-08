package system

import (
	"context"
	"io"
)

// Environment represents a execution environment.
type Environment interface {
	Contextual
	Outputs
}

// Contextual returns a context.Context object.
type Contextual interface {
	Context() context.Context
}

// Outputs holds current program outputs.
type Outputs interface {
	OutOrStdout() io.Writer
	ErrOrStderr() io.Writer
}

// WithOutputs returns a new Environment with the given outputs.
func WithOutputs(out, err io.Writer, env Environment) Environment {
	return &outputs{out, err, env}
}

// WithContext returns a new Environment with the given context.
func WithContext(ctx context.Context, env Environment) Environment {
	return &contextual{ctx, env}
}

type outputs struct {
	out, err io.Writer
	env      Environment
}

func (o outputs) Context() context.Context {
	return o.env.Context()
}

func (o outputs) OutOrStdout() io.Writer {
	return o.out
}

func (o outputs) ErrOrStderr() io.Writer {
	return o.err
}

type contextual struct {
	ctx context.Context
	env Environment
}

func (c contextual) Context() context.Context {
	return c.ctx
}

func (c contextual) OutOrStdout() io.Writer {
	return c.env.OutOrStdout()
}

func (c contextual) ErrOrStderr() io.Writer {
	return c.env.ErrOrStderr()
}
