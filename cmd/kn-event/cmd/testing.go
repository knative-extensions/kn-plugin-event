package cmd

import "io"

// TestingCmd a wrapper for Cmd to ease of testing.
type TestingCmd struct {
	*Cmd
}

// ExecuteOrFail the command or fail on error.
func (c *TestingCmd) ExecuteOrFail() {
	c.init()
	c.Cmd.Execute()
}

// Execute the command and return error if any.
func (c *TestingCmd) Execute() error {
	c.init()
	return c.execute()
}

// Exit sets the exit command that accepts retcode.
func (c *TestingCmd) Exit(fn func(code int)) {
	c.init()
	c.exit = fn
}

// Out sets output stream to cmd.
func (c *TestingCmd) Out(newOut io.Writer) {
	c.init()
	c.root.SetOut(newOut)
	c.root.SetErr(newOut)
}

// Args set to main command to be executed.
func (c *TestingCmd) Args(args ...string) {
	c.init()
	c.root.SetArgs(args)
}

func (c *TestingCmd) init() {
	if c.Cmd == nil {
		c.Cmd = &Cmd{}
	}
	c.Cmd.init()
}
