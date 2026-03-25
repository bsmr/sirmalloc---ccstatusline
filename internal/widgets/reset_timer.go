package widgets

import (
	"fmt"
	"time"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/util"
)

func init() { Register("reset-timer", resetTimerWidget{}) }

type resetTimerWidget struct{}

func (resetTimerWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil || si.TranscriptPath == "" {
		return ""
	}
	blockStart, err := getBlockStart(si)
	if err != nil {
		return ""
	}

	const fiveHours = 5 * time.Hour
	elapsed := time.Since(blockStart)
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed > fiveHours {
		elapsed = fiveHours
	}
	remaining := fiveHours - elapsed

	displayMode := item.Metadata["display"]
	switch displayMode {
	case "progress":
		pct := float64(elapsed) / float64(fiveHours) * 100
		bar := makeProgressBar(pct, 32)
		label := fmt.Sprintf("[%s] %.1f%%", bar, pct)
		if item.RawValue {
			return label
		}
		return "Reset " + label
	case "progress-short":
		pct := float64(elapsed) / float64(fiveHours) * 100
		bar := makeProgressBar(pct, 16)
		label := fmt.Sprintf("[%s] %.1f%%", bar, pct)
		if item.RawValue {
			return label
		}
		return "Reset " + label
	default:
		label := util.FormatDuration(remaining.Milliseconds())
		if item.RawValue {
			return label
		}
		return "Reset: " + label
	}
}
