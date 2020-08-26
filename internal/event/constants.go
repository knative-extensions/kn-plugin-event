package event

import (
	"fmt"

	"github.com/cardil/kn-event/internal"
	"github.com/google/uuid"
)

const (
	// DefaultType holds a default type for a event
	DefaultType = "dev.knative.cli.plugin.event.generic"
)

var (
	// DefaultSource holds a default source of an event
	DefaultSource = fmt.Sprintf("%s/%s", internal.PluginName, internal.Version)
)

// NewID creates a new ID for an event
func NewID() string {
	return uuid.New().String()
}
