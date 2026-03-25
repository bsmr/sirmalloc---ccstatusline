package widgets

import (
	"fmt"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("context-percentage", contextPctWidget{})
}

type contextPctWidget struct{}

func (contextPctWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil {
		return ""
	}
	pct := si.ContextWindow.UsedPercentage
	if item.Remaining {
		pct = si.ContextWindow.RemainingPct
	}
	return fmt.Sprintf("%.4g%%", pct)
}
