package widgets

import (
	"fmt"
	"strings"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() { Register("context-bar", contextBarWidget{}) }

type contextBarWidget struct{}

func (contextBarWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil {
		return ""
	}
	cw := si.ContextWindow
	total := cw.ContextWindowSize
	used := cw.TotalInputTokens + cw.TotalOutputTokens
	if total <= 0 {
		return ""
	}

	barWidth := 16
	if item.Metadata["display"] == "progress" {
		barWidth = 32
	}

	pct := float64(used) / float64(total) * 100
	if pct > 100 {
		pct = 100
	}
	if pct < 0 {
		pct = 0
	}

	bar := makeProgressBar(pct, barWidth)
	usedK := (used + 500) / 1000
	totalK := (total + 500) / 1000
	display := fmt.Sprintf("%s %dk/%dk (%d%%)", bar, usedK, totalK, int(pct))

	if item.RawValue {
		return display
	}
	return "Context: " + display
}

// makeProgressBar renders a progress bar of the given width for a percentage value.
// It is shared across context_bar, reset_timer, and session_usage widgets.
func makeProgressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
