package config

import "github.com/fatih/color"

var (
	// RepoDir holds a repository path.
	RepoDir string

	// BuildDirPath holds a build dir path.
	BuildDirPath = []string{"build", "_output"}
	// MageTag holds default mage tag settings.
	MageTag = MageTagStruct{
		Color: color.FgCyan,
		Label: "[MAGE]",
	}

	// Dependencies will hold additional dependencies that needs to be installed
	// before running tasks.
	Dependencies = []string{
		"github.com/kyoh86/richgo",
	}

	// VersionVariablePath a Golang path to version holding variable.
	VersionVariablePath string

	// Binaries a list of binaries to build.
	Binaries []Binary

	// CleaningTasks additional cleaning tasks.
	CleaningTasks []CustomTask

	// Checks holds a list of checks to perform.
	Checks []CustomTask
)
