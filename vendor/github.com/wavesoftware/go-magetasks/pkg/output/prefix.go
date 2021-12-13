package output

import (
	"github.com/fatih/color"
	"github.com/wavesoftware/go-magetasks/config"
)

func prefix() string {
	mt := config.Actual().MageTag
	return color.New(mt.Color).Sprint(mt.Label)
}
