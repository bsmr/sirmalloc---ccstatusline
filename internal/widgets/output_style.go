package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("output-style", outputStyleWidget{})
}

type outputStyleWidget struct{}

func (outputStyleWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.OutputStyle == nil {
		return ""
	}
	return si.OutputStyle.Name
}
