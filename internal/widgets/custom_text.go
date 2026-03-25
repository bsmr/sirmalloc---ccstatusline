package widgets

import (
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("custom-text", customTextWidget{})
}

type customTextWidget struct{}

func (customTextWidget) Render(_ *input.StatusInput, item config.WidgetItem) string {
	return item.Text
}
