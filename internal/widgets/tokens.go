package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/util"
)

func init() {
	Register("tokens-input", tokensInputWidget{})
	Register("tokens-output", tokensOutputWidget{})
	Register("tokens-cached", tokensCachedWidget{})
	Register("tokens-total", tokensTotalWidget{})
}

type tokensInputWidget struct{}

func (tokensInputWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil {
		return ""
	}
	return util.FormatTokens(si.ContextWindow.TotalInputTokens)
}

type tokensOutputWidget struct{}

func (tokensOutputWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil {
		return ""
	}
	return util.FormatTokens(si.ContextWindow.TotalOutputTokens)
}

type tokensCachedWidget struct{}

func (tokensCachedWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil || si.ContextWindow.CurrentUsage == nil {
		return ""
	}
	return util.FormatTokens(si.ContextWindow.CurrentUsage.CacheReadInputTokens)
}

type tokensTotalWidget struct{}

func (tokensTotalWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.ContextWindow == nil {
		return ""
	}
	total := si.ContextWindow.TotalInputTokens + si.ContextWindow.TotalOutputTokens
	return util.FormatTokens(total)
}
