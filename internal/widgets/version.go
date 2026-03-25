package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("version", versionWidget{})
}

type versionWidget struct{}

func (versionWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil {
		return ""
	}
	return si.Version
}
