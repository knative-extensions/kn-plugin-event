package config

import "github.com/fatih/color"

// MageTagStruct holds a mage tag.
type MageTagStruct struct {
	Color color.Attribute
	Label string
}

// Binary represents a binary that will be built.
type Binary struct {
	Name      string
	ImageArgs map[string]string
}

// CustomTask is a custom function that will be used in the build.
type CustomTask struct {
	Name string
	Task func() error
}
