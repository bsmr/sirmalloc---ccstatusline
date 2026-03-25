package render

import (
	"strings"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/widgets"
)

const linePad = " "

// widgetSlot holds the rendered output for a single widget position in a line.
type widgetSlot struct {
	value  string           // rendered text (may include ANSI)
	isFlex bool             // true when this is a flex-separator placeholder
	item   config.WidgetItem
}

// Render converts a StatusInput + Settings into one rendered string per line.
// termWidth is used for flex-separator layout and final truncation.
func Render(si *input.StatusInput, cfg *config.Settings, termWidth int) []string {
	width := effectiveWidth(termWidth, cfg, si)

	lines := make([]string, 0, len(cfg.Lines))
	for _, lineItems := range cfg.Lines {
		lines = append(lines, renderLine(si, cfg, lineItems, width))
	}
	return lines
}

// renderLine renders a single line of widgets.
func renderLine(si *input.StatusInput, cfg *config.Settings, lineItems []config.WidgetItem, width int) string {
	// --- Pass 1: collect slots ---
	slots := make([]widgetSlot, 0, len(lineItems))
	totalVisible := 0
	flexCount := 0
	nonFlexCount := 0

	for _, item := range lineItems {
		w, ok := widgets.Get(item.Type)
		if !ok {
			continue
		}

		if item.Type == "flex-separator" {
			slots = append(slots, widgetSlot{isFlex: true, item: item})
			flexCount++
			continue
		}

		val := w.Render(si, item)
		if val == "" {
			continue
		}

		val = applyColors(val, item, cfg)
		slots = append(slots, widgetSlot{value: val, item: item})
		totalVisible += VisibleLen(val)
		nonFlexCount++
	}

	if nonFlexCount == 0 && flexCount == 0 {
		return ""
	}

	// --- Flex width calculation ---
	// Layout (standard): <pad> slot <pad> slot <pad> ...
	// Fixed pads = (nonFlexCount + flexCount + 1) × len(linePad)
	if flexCount > 0 {
		fixedWidth := totalVisible + (nonFlexCount+flexCount+1)*len(linePad)
		available := width - fixedWidth
		if available < 0 {
			available = 0
		}
		flexWidth := available / flexCount

		// Resolve flex slots to their space-fill values.
		for i := range slots {
			if slots[i].isFlex {
				slots[i].value = strings.Repeat(" ", flexWidth)
				slots[i].isFlex = false
			}
		}
	}

	// --- Powerline mode ---
	if cfg.Powerline != nil && cfg.Powerline.Enabled {
		segs := slotsToSegments(slots, cfg)
		if len(segs) == 0 {
			return ""
		}
		return renderPowerline(segs, cfg.Powerline)
	}

	// --- Standard join ---
	parts := make([]string, 0, len(slots))
	for _, s := range slots {
		parts = append(parts, s.value)
	}
	line := linePad + strings.Join(parts, linePad) + linePad
	line = strings.TrimRight(line, " ")
	return Truncate(line, width)
}

// slotsToSegments converts widgetSlots into powerline segments,
// skipping slots with empty values.
func slotsToSegments(slots []widgetSlot, cfg *config.Settings) []segment {
	segs := make([]segment, 0, len(slots))
	for _, s := range slots {
		if s.value == "" {
			continue
		}
		segs = append(segs, segment{
			text: s.value,
			fg:   s.item.FG,
			bg:   s.item.BG,
			bold: s.item.Bold || cfg.GlobalBold,
		})
	}
	return segs
}

// applyColors wraps val with ANSI color/bold codes derived from item and global
// settings. Returns val unchanged when no styling is configured.
func applyColors(val string, item config.WidgetItem, cfg *config.Settings) string {
	prefix := FG(item.FG) + BG(item.BG)
	if item.Bold || cfg.GlobalBold {
		prefix += Bold()
	}
	if prefix == "" {
		return val
	}
	return prefix + val + Reset()
}

// effectiveWidth returns the usable terminal width according to cfg.FlexMode.
func effectiveWidth(termWidth int, cfg *config.Settings, si *input.StatusInput) int {
	switch cfg.FlexMode {
	case "full":
		return termWidth
	case "full-until-compact":
		threshold := cfg.CompactThreshold
		if threshold <= 0 {
			threshold = 80
		}
		if si != nil && si.ContextWindow != nil &&
			si.ContextWindow.UsedPercentage >= float64(threshold) {
			return termWidth / 2
		}
		return termWidth
	default: // "full-minus-40"
		w := termWidth - 40
		if w < 40 {
			w = 40
		}
		return w
	}
}
