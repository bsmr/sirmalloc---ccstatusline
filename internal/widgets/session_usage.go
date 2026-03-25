package widgets

import (
	"fmt"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() { Register("session-usage", sessionUsageWidget{}) }

type sessionUsageWidget struct{}

func (sessionUsageWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil || si.RateLimits == nil || si.RateLimits.FiveHour == nil || si.RateLimits.FiveHour.UsedPercentage == nil {
		return ""
	}
	pct := *si.RateLimits.FiveHour.UsedPercentage
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}

	displayMode := item.Metadata["display"]
	var label string
	switch displayMode {
	case "progress":
		bar := makeProgressBar(pct, 32)
		label = fmt.Sprintf("%s %.1f%%", bar, pct)
	case "progress-short":
		bar := makeProgressBar(pct, 16)
		label = fmt.Sprintf("%s %.1f%%", bar, pct)
	default:
		label = fmt.Sprintf("%.1f%%", pct)
	}

	if item.RawValue {
		return label
	}
	return "Session: " + label
}
