package widgets

import (
	"fmt"
	"strings"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/util"
)

func init() {
	Register("context-length", contextLengthWidget{})
	Register("context-pct-usable", contextPctUsableWidget{})
}

type contextLengthWidget struct{}

func (contextLengthWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil {
		return ""
	}
	return util.FormatTokens(si.ContextWindow.ContextWindowSize)
}

type contextPctUsableWidget struct{}

func getModelContext(modelID string) (maxTokens, usableTokens int) {
	if strings.Contains(strings.ToLower(modelID), "1m") {
		return 1_000_000, 800_000
	}
	return 200_000, 160_000
}

func (contextPctUsableWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil {
		return ""
	}
	modelID := ""
	if si.Model != nil {
		modelID = si.Model.ID
	}
	_, usableTokens := getModelContext(modelID)
	usedTokens := si.ContextWindow.TotalInputTokens + si.ContextWindow.TotalOutputTokens
	pct := float64(usableTokens-usedTokens) / float64(usableTokens) * 100
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	return fmt.Sprintf("%.0f%%", pct)
}
