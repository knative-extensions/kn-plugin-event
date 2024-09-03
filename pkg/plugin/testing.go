package plugin

import (
	"bytes"
	"io"

	knplugin "knative.dev/client/pkg/plugin"
)

// WithCapture captures the output from of a running plugin.
func WithCapture(fn func()) []byte {
	var buf bytes.Buffer
	WithOutput(&buf, fn)
	return buf.Bytes()
}

// WithOutput executes provided function after configuring the stdout and stderr.
func WithOutput(output io.Writer, fn func()) {
	pl := findPlugin()
	save := pl.Writer
	pl.Writer = output
	defer func() {
		pl.Writer = save
	}()
	fn()
}

func findPlugin() *plugin {
	for _, ip := range knplugin.InternalPlugins {
		if pl, ok := ip.(*plugin); ok {
			return pl
		}
	}
	panic("this should never happen")
}
