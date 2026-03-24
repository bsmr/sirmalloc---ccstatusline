package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("model", modelWidget{})
}

type modelWidget struct{}

func (modelWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si.Model == nil {
		return ""
	}
	name := si.Model.DisplayName
	if name == "" {
		name = si.Model.ID
	}
	if name == "" {
		return ""
	}
	if item.RawValue {
		return name
	}
	return "Model: " + name
}
