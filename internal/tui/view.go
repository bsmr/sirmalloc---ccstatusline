package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go.a8l.eu/ccstatusline/internal/config"
)

// lipgloss color codes for the color picker (ANSI 0–7, empty = default/reset).
var colorLipgloss = []lipgloss.Color{
	lipgloss.Color("0"), // black
	lipgloss.Color("1"), // red
	lipgloss.Color("2"), // green
	lipgloss.Color("3"), // yellow
	lipgloss.Color("4"), // blue
	lipgloss.Color("5"), // magenta
	lipgloss.Color("6"), // cyan
	lipgloss.Color("7"), // white
	lipgloss.Color(""),  // default
}

var (
	styleSelected = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleDim      = lipgloss.NewStyle().Faint(true)
	styleErr      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func (m model) View() string {
	var b strings.Builder
	switch m.screen {
	case screenLines:
		m.viewLines(&b)
	case screenLine:
		m.viewLine(&b)
	case screenWidget:
		m.viewWidget(&b)
	case screenTypeSelect:
		m.viewTypeSelect(&b)
	case screenGlobalOptions:
		m.viewGlobalOptions(&b)
	case screenPowerline:
		m.viewPowerline(&b)
	case screenInstall:
		m.viewInstall(&b)
	case screenPreview:
		m.viewPreview(&b)
	case screenColorPicker:
		m.viewColorPicker(&b)
	}
	b.WriteString("\n")
	if m.statusMsg != "" {
		b.WriteString(styleErr.Render(m.statusMsg))
	} else {
		b.WriteString(styleDim.Render(m.hints()))
	}
	b.WriteString("\n")
	return b.String()
}

func (m model) hints() string {
	switch m.screen {
	case screenLines:
		return "j/k navigate · Enter open · a add · d delete · g global opts · p powerline · v preview · i install · Ctrl+S save · q quit"
	case screenLine:
		return "j/k navigate · Enter edit · a add · d delete · Ctrl+S save · Esc back"
	case screenWidget:
		if m.editingText {
			return "type · Enter/Esc confirm"
		}
		return "j/k navigate · Enter/Space toggle · c color picker · Ctrl+S save+back · Esc back"
	case screenTypeSelect:
		return "j/k navigate · Enter select · Esc cancel"
	case screenGlobalOptions:
		if m.editingText {
			return "type · Enter/Esc confirm"
		}
		return "j/k navigate · Enter/Space toggle · Ctrl+S save+back · Esc/q back"
	case screenPowerline:
		if m.editingText {
			return "type · Enter/Esc confirm"
		}
		return "j/k navigate · Enter/Space toggle · Ctrl+S save+back · Esc/q back"
	case screenInstall:
		return "i install · u uninstall · Esc back"
	case screenPreview:
		return "r refresh · Esc back"
	case screenColorPicker:
		return "j/k navigate · Enter select · Esc cancel"
	}
	return ""
}

func (m model) viewLines(b *strings.Builder) {
	b.WriteString("Lines\n\n")
	for i, line := range m.settings.Lines {
		label := fmt.Sprintf("Line %d  (%d widgets)", i+1, len(line))
		if i == m.linesCursor {
			b.WriteString(styleSelected.Render("> "+label) + "\n")
		} else {
			b.WriteString("  " + label + "\n")
		}
	}
	if len(m.settings.Lines) == 0 {
		b.WriteString(styleDim.Render("  (no lines — press 'a' to add one)") + "\n")
	}
}

func (m model) viewLine(b *strings.Builder) {
	b.WriteString(fmt.Sprintf("Line %d\n\n", m.linesCursor+1))
	line := m.settings.Lines[m.linesCursor]
	for i, w := range line {
		summary := widgetSummary(w)
		label := fmt.Sprintf("[%d] %-22s %s", i+1, w.Type, summary)
		if i == m.widgetsCursor {
			b.WriteString(styleSelected.Render("> "+label) + "\n")
		} else {
			b.WriteString("  " + label + "\n")
		}
	}
	if len(line) == 0 {
		b.WriteString(styleDim.Render("  (empty — press 'a' to add a widget)") + "\n")
	}
}

func widgetSummary(w config.WidgetItem) string {
	var parts []string
	if w.FG != "" {
		parts = append(parts, "fg="+w.FG)
	}
	if w.BG != "" {
		parts = append(parts, "bg="+w.BG)
	}
	if w.Bold {
		parts = append(parts, "bold")
	}
	switch w.Type {
	case "separator":
		sep := w.SepChar
		if sep == "" {
			sep = "|"
		}
		parts = append(parts, "sep="+sep)
	case "context-percentage":
		if w.Remaining {
			parts = append(parts, "remaining")
		}
	case "cwd":
		if w.Segments > 0 {
			parts = append(parts, fmt.Sprintf("segments=%d", w.Segments))
		}
		if w.FishStyle {
			parts = append(parts, "fishStyle")
		}
	case "custom-text":
		if w.Text != "" {
			parts = append(parts, "text="+w.Text)
		}
	case "custom-command":
		if w.Command != "" {
			parts = append(parts, "cmd="+w.Command)
		}
	case "git-branch", "git-changes", "git-worktree":
		if w.HideNoGit {
			parts = append(parts, "hideNoGit")
		}
	case "block-timer":
		if w.BarMode != "" {
			parts = append(parts, "barMode="+w.BarMode)
		}
	}
	return strings.Join(parts, " ")
}

func (m model) viewWidget(b *strings.Builder) {
	w := m.currentWidget()
	if w == nil {
		return
	}
	fields := fieldsFor(w)
	b.WriteString(fmt.Sprintf("Widget: %s\n\n", w.Type))
	for i, f := range fields {
		val := getField(w, i)
		var label string
		if i == m.fieldCursor && m.editingText {
			label = fmt.Sprintf("%-14s  %s_", f.label, m.textBuf)
		} else {
			label = fmt.Sprintf("%-14s  %s", f.label, val)
		}
		if i == m.fieldCursor {
			b.WriteString(styleSelected.Render("> "+label) + "\n")
		} else {
			b.WriteString("  " + label + "\n")
		}
	}
}

func (m model) viewTypeSelect(b *strings.Builder) {
	b.WriteString("Add widget — select type\n\n")
	for i, t := range allWidgetTypes {
		if i == m.typeCursor {
			b.WriteString(styleSelected.Render("> "+t) + "\n")
		} else {
			b.WriteString("  " + t + "\n")
		}
	}
}

func (m model) viewGlobalOptions(b *strings.Builder) {
	b.WriteString("Global Options\n\n")
	for i, label := range globalFieldLabels {
		val := getGlobalField(m.settings, i)
		var entry string
		if i == m.fieldCursor && m.editingText {
			entry = fmt.Sprintf("%-24s  %s_", label, m.textBuf)
		} else {
			entry = fmt.Sprintf("%-24s  %s", label, val)
		}
		if i == m.fieldCursor {
			b.WriteString(styleSelected.Render("> "+entry) + "\n")
		} else {
			b.WriteString("  " + entry + "\n")
		}
	}
}

func (m model) viewPowerline(b *strings.Builder) {
	b.WriteString("Powerline\n\n")
	for i, label := range powerlineFieldLabels {
		val := getPowerlineField(m.settings, i)
		var entry string
		if i == m.fieldCursor && m.editingText {
			entry = fmt.Sprintf("%-12s  %s_", label, m.textBuf)
		} else {
			entry = fmt.Sprintf("%-12s  %s", label, val)
		}
		if i == m.fieldCursor {
			b.WriteString(styleSelected.Render("> "+entry) + "\n")
		} else {
			b.WriteString("  " + entry + "\n")
		}
	}
}

// ---------------------------------------------------------------------------
// Install screen
// ---------------------------------------------------------------------------

func (m model) viewInstall(b *strings.Builder) {
	b.WriteString("Install / Uninstall\n\n")

	installed, cmd, err := config.IsInstalled()
	switch {
	case err != nil:
		b.WriteString(styleErr.Render("error checking status: "+err.Error()) + "\n")
	case installed:
		b.WriteString(fmt.Sprintf("Current status: %s  (%s)\n",
			styleSelected.Render("installed"), cmd))
	default:
		b.WriteString("Current status: not installed\n")
	}

	b.WriteString("\n")
	b.WriteString("  i  install in Claude Code settings\n")
	b.WriteString("  u  uninstall from Claude Code settings\n")
	b.WriteString("\n")
	b.WriteString(styleDim.Render("Esc  back") + "\n")
}

// ---------------------------------------------------------------------------
// Preview screen
// ---------------------------------------------------------------------------

func (m model) viewPreview(b *strings.Builder) {
	b.WriteString("Preview  (r to refresh)\n\n")
	if len(m.previewLines) == 0 {
		b.WriteString(styleDim.Render("  (no output)") + "\n")
		return
	}
	for _, line := range m.previewLines {
		b.WriteString("  " + line + "\n")
	}
}

// ---------------------------------------------------------------------------
// Color picker screen
// ---------------------------------------------------------------------------

func (m model) viewColorPicker(b *strings.Builder) {
	b.WriteString(fmt.Sprintf("Pick color for %s\n\n", m.colorPickerTarget))
	for i, opt := range colorOptions {
		var square string
		if opt.code == "" {
			square = styleDim.Render("█")
		} else {
			square = lipgloss.NewStyle().Foreground(colorLipgloss[i]).Render("█")
		}
		label := fmt.Sprintf("%s %s", square, opt.name)
		if i == m.colorCursor {
			b.WriteString(styleSelected.Render("> "+label) + "\n")
		} else {
			b.WriteString("  " + label + "\n")
		}
	}
}
