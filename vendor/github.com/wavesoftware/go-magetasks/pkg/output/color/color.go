package color

import (
	"os"

	"github.com/fatih/color"
)

var (
	// Red text color.
	Red = color.New(color.FgHiRed).Add(color.Bold).SprintFunc()
	// Green text color.
	Green = color.New(color.FgHiGreen).Add(color.Bold).SprintFunc()
	// Yellow text color.
	Yellow = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	// Blue text color.
	Blue = color.New(color.FgBlue).SprintFunc()
)

// SetupMode will set the output color mode.
func SetupMode() {
	if val, envset := os.LookupEnv("FORCE_COLOR"); envset && val == "true" {
		color.NoColor = false
	}
}
