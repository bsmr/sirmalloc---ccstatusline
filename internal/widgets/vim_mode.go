package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("vim-mode", vimModeWidget{})
}

type vimModeWidget struct{}

func (vimModeWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.Vim == nil {
		return ""
	}
	return si.Vim.Mode
}
