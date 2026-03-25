package widgets

import (
	"fmt"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

var currentTermWidth int

// SetTermWidth stores the terminal width for use by the terminal-width widget.
func SetTermWidth(w int) { currentTermWidth = w }

func init() {
	Register("terminal-width", terminalWidthWidget{})
}

type terminalWidthWidget struct{}

func (terminalWidthWidget) Render(_ *input.StatusInput, _ config.WidgetItem) string {
	return fmt.Sprintf("%d", currentTermWidth)
}
