package tasks

import (
	"github.com/fatih/color"
	"github.com/wavesoftware/go-magetasks/config"
)

func mageTag() string {
	return color.New(config.MageTag.Color).Sprint(config.MageTag.Label)
}
