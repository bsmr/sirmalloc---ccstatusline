package render

import (
	"strings"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/widgets"
)

// Render converts a StatusInput + Settings into one rendered string per line.
// termWidth is reserved for future truncation (Phase 3) and currently ignored.
func Render(si *input.StatusInput, cfg *config.Settings, _ int) []string {
	const pad = " "

	lines := make([]string, 0, len(cfg.Lines))
	for _, lineItems := range cfg.Lines {
		var parts []string
		for _, item := range lineItems {
			w, ok := widgets.Get(item.Type)
			if !ok {
				continue
			}
			val := w.Render(si, item)
			if val == "" {
				continue
			}
			parts = append(parts, val)
		}
		if len(parts) == 0 {
			lines = append(lines, "")
			continue
		}
		// Join parts with single-space padding, trim trailing space.
		line := pad + strings.Join(parts, pad) + pad
		lines = append(lines, strings.TrimRight(line, " "))
	}
	return lines
}
