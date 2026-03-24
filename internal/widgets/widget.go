package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

// Widget is the interface every status-line widget must implement.
type Widget interface {
	Render(si *input.StatusInput, item config.WidgetItem) string
}

var registry = map[string]Widget{}

// Register adds a widget implementation under the given type name.
func Register(typ string, w Widget) {
	registry[typ] = w
}

// Get returns the widget registered under typ, or false if unknown.
func Get(typ string) (Widget, bool) {
	w, ok := registry[typ]
	return w, ok
}
