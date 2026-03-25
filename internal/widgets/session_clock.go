package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/util"
)

func init() {
	Register("session-clock", sessionClockWidget{})
}

type sessionClockWidget struct{}

func (sessionClockWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.Cost == nil {
		return ""
	}
	return util.FormatDuration(si.Cost.TotalDurationMs)
}
