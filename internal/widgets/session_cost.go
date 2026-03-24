package widgets

import (
	"fmt"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("session-cost", sessionCostWidget{})
}

type sessionCostWidget struct{}

func (sessionCostWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si.Cost == nil {
		return ""
	}
	return fmt.Sprintf("$%.4f", si.Cost.TotalCostUSD)
}
