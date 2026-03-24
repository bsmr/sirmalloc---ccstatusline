package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("flex-separator", flexSeparatorWidget{})
}

type flexSeparatorWidget struct{}

// Render returns an empty string; flex layout is handled by the renderer in Phase 3.
func (flexSeparatorWidget) Render(_ *input.StatusInput, _ config.WidgetItem) string {
	return ""
}
