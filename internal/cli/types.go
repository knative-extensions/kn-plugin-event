package cli

import "github.com/thediveo/enumflag"

// EventArgs holds args of event to be created with
type EventArgs struct {
	Type      string
	ID        string
	Source    string
	Fields    []string
	RawFields []string
}

// OutputMode is type of output to produce
type OutputMode enumflag.Flag

// OutputMode enumeration values.
const (
	HumanReadable OutputMode = iota
	JSON
	YAML
)
