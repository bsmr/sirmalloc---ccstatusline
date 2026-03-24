package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("separator", separatorWidget{})
}

type separatorWidget struct{}

func (separatorWidget) Render(_ *input.StatusInput, item config.WidgetItem) string {
	if item.SepChar != "" {
		return item.SepChar
	}
	return "|"
}
