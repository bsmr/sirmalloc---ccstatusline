package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/render"
)

type screenKind int

const (
	screenLines      screenKind = iota
	screenLine
	screenWidget
	screenTypeSelect
	screenGlobalOptions
	screenPowerline
	screenInstall
	screenPreview
	screenColorPicker
)

// allWidgetTypes lists all widget types available for addition.
var allWidgetTypes = []string{
	"model",
	"context-percentage",
	"context-pct-usable",
	"context-length",
	"context-bar",
	"reset-timer",
	"session-usage",
	"session-clock",
	"session-cost",
	"block-timer",
	"git-branch",
	"git-changes",
	"git-worktree",
	"cwd",
	"version",
	"output-style",
	"tokens-input",
	"tokens-output",
	"tokens-cached",
	"tokens-total",
	"rate-limit-five-hour",
	"rate-limit-seven-day",
	"vim-mode",
	"terminal-width",
	"custom-text",
	"custom-command",
	"separator",
	"flex-separator",
}

type widgetField struct {
	label  string
	isText bool // false = bool toggle
}

func fieldsFor(w *config.WidgetItem) []widgetField {
	fields := []widgetField{
		{"fg", true},
		{"bg", true},
		{"bold", false},
		{"backgroundColor", true},
	}
	switch w.Type {
	case "separator":
		fields = append(fields, widgetField{"sepChar", true})
	case "context-percentage":
		fields = append(fields, widgetField{"remaining", false})
	case "cwd":
		fields = append(fields, widgetField{"segments", true})
		fields = append(fields, widgetField{"fishStyle", false})
	case "custom-text":
		fields = append(fields, widgetField{"text", true})
	case "custom-command":
		fields = append(fields, widgetField{"command", true})
		fields = append(fields, widgetField{"timeout", true})
	case "git-branch", "git-changes", "git-worktree":
		fields = append(fields, widgetField{"hideNoGit", false})
	case "block-timer":
		fields = append(fields, widgetField{"barMode", true})
	case "context-bar", "session-usage", "reset-timer":
		// display mode is in metadata["display"] — edit as text
		fields = append(fields, widgetField{"metadata.display", true})
	}
	return fields
}

func getField(w *config.WidgetItem, idx int) string {
	fields := fieldsFor(w)
	if idx >= len(fields) {
		return ""
	}
	switch fields[idx].label {
	case "fg":
		return w.FG
	case "bg":
		return w.BG
	case "bold":
		if w.Bold {
			return "true"
		}
		return "false"
	case "sepChar":
		if w.SepChar == "" {
			return "|"
		}
		return w.SepChar
	case "remaining":
		if w.Remaining {
			return "true"
		}
		return "false"
	case "segments":
		return strconv.Itoa(w.Segments)
	case "fishStyle":
		if w.FishStyle {
			return "true"
		}
		return "false"
	case "hideNoGit":
		if w.HideNoGit {
			return "true"
		}
		return "false"
	case "text":
		return w.Text
	case "command":
		return w.Command
	case "timeout":
		return strconv.Itoa(w.Timeout)
	case "barMode":
		return w.BarMode
	case "backgroundColor":
		return w.BackgroundColor
	case "metadata.display":
		if w.Metadata == nil {
			return ""
		}
		return w.Metadata["display"]
	}
	return ""
}

func setField(w *config.WidgetItem, idx int, val string) {
	fields := fieldsFor(w)
	if idx >= len(fields) {
		return
	}
	switch fields[idx].label {
	case "fg":
		w.FG = val
	case "bg":
		w.BG = val
	case "bold":
		w.Bold = val == "true"
	case "sepChar":
		w.SepChar = val
	case "remaining":
		w.Remaining = val == "true"
	case "segments":
		if n, err := strconv.Atoi(val); err == nil {
			w.Segments = n
		}
	case "fishStyle":
		w.FishStyle = val == "true"
	case "hideNoGit":
		w.HideNoGit = val == "true"
	case "text":
		w.Text = val
	case "command":
		w.Command = val
	case "timeout":
		if n, err := strconv.Atoi(val); err == nil {
			w.Timeout = n
		}
	case "barMode":
		w.BarMode = val
	case "backgroundColor":
		w.BackgroundColor = val
	case "metadata.display":
		if w.Metadata == nil {
			w.Metadata = make(map[string]string)
		}
		w.Metadata["display"] = val
	}
}

func toggleBool(w *config.WidgetItem, idx int) {
	fields := fieldsFor(w)
	if idx >= len(fields) || fields[idx].isText {
		return
	}
	switch fields[idx].label {
	case "bold":
		w.Bold = !w.Bold
	case "remaining":
		w.Remaining = !w.Remaining
	case "fishStyle":
		w.FishStyle = !w.FishStyle
	case "hideNoGit":
		w.HideNoGit = !w.HideNoGit
	}
}

// colorOptions maps cursor index to a display name and the ANSI/lipgloss color code.
var colorOptions = []struct{ name, code string }{
	{"black", "black"},
	{"red", "red"},
	{"green", "green"},
	{"yellow", "yellow"},
	{"blue", "blue"},
	{"magenta", "magenta"},
	{"cyan", "cyan"},
	{"white", "white"},
	{"default", ""},
}

type model struct {
	settings   *config.Settings
	configPath string
	screen     screenKind

	linesCursor   int
	widgetsCursor int
	fieldCursor   int
	typeCursor    int

	editingText bool
	textBuf     string

	statusMsg string

	// Install screen
	installStatus string

	// Preview screen
	previewLines []string

	// Color picker screen
	colorPickerTarget string // "fg" or "bg"
	colorCursor       int
}

func newModel(s *config.Settings, path string) model {
	return model{settings: s, configPath: path}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	m.statusMsg = ""
	switch m.screen {
	case screenLines:
		return m.updateLines(key)
	case screenLine:
		return m.updateLine(key)
	case screenWidget:
		return m.updateWidget(key)
	case screenTypeSelect:
		return m.updateTypeSelect(key)
	case screenGlobalOptions:
		return m.updateGlobalOptions(key)
	case screenPowerline:
		return m.updatePowerline(key)
	case screenInstall:
		return m.updateInstall(key)
	case screenPreview:
		return m.updatePreview(key)
	case screenColorPicker:
		return m.updateColorPicker(key)
	}
	return m, nil
}

func (m model) clampLines() model {
	n := len(m.settings.Lines)
	if n == 0 {
		m.linesCursor = 0
		return m
	}
	if m.linesCursor >= n {
		m.linesCursor = n - 1
	}
	if m.linesCursor < 0 {
		m.linesCursor = 0
	}
	return m
}

func (m model) clampWidgets() model {
	if m.linesCursor >= len(m.settings.Lines) {
		return m
	}
	n := len(m.settings.Lines[m.linesCursor])
	if n == 0 {
		m.widgetsCursor = 0
		return m
	}
	if m.widgetsCursor >= n {
		m.widgetsCursor = n - 1
	}
	if m.widgetsCursor < 0 {
		m.widgetsCursor = 0
	}
	return m
}

func (m model) save() model {
	if err := config.Save(m.configPath, m.settings); err != nil {
		m.statusMsg = "save error: " + err.Error()
	} else {
		m.statusMsg = "saved."
	}
	return m
}

func (m model) updateLines(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "j", "down":
		m.linesCursor++
		m = m.clampLines()
	case "k", "up":
		m.linesCursor--
		m = m.clampLines()
	case "enter":
		if len(m.settings.Lines) > 0 {
			m.widgetsCursor = 0
			m.screen = screenLine
		}
	case "a":
		m.settings.Lines = append(m.settings.Lines, []config.WidgetItem{})
		m.linesCursor = len(m.settings.Lines) - 1
	case "d":
		if len(m.settings.Lines) > 1 {
			m.settings.Lines = append(
				m.settings.Lines[:m.linesCursor],
				m.settings.Lines[m.linesCursor+1:]...,
			)
			m = m.clampLines()
		} else {
			m.statusMsg = "cannot delete the last line"
		}
	case "g":
		m.screen = screenGlobalOptions
		m.fieldCursor = 0
	case "p":
		m.screen = screenPowerline
		m.fieldCursor = 0
	case "v":
		m.previewLines = m.renderPreview()
		m.screen = screenPreview
	case "i":
		m.installStatus = ""
		m.screen = screenInstall
	case "ctrl+s":
		m = m.save()
	}
	return m, nil
}

func (m model) updateLine(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenLines
	case "j", "down":
		m.widgetsCursor++
		m = m.clampWidgets()
	case "k", "up":
		m.widgetsCursor--
		m = m.clampWidgets()
	case "enter":
		line := m.settings.Lines[m.linesCursor]
		if len(line) > 0 {
			m.fieldCursor = 0
			m.editingText = false
			m.screen = screenWidget
		}
	case "a":
		m.typeCursor = 0
		m.screen = screenTypeSelect
	case "d":
		line := m.settings.Lines[m.linesCursor]
		if len(line) > 0 {
			m.settings.Lines[m.linesCursor] = append(
				line[:m.widgetsCursor],
				line[m.widgetsCursor+1:]...,
			)
			m = m.clampWidgets()
		}
	case "ctrl+s":
		m = m.save()
	}
	return m, nil
}

func (m model) currentWidget() *config.WidgetItem {
	if m.linesCursor >= len(m.settings.Lines) {
		return nil
	}
	line := m.settings.Lines[m.linesCursor]
	if m.widgetsCursor >= len(line) {
		return nil
	}
	return &m.settings.Lines[m.linesCursor][m.widgetsCursor]
}

func (m model) updateWidget(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	w := m.currentWidget()
	if w == nil {
		m.screen = screenLine
		return m, nil
	}
	fields := fieldsFor(w)

	if m.editingText {
		switch key.String() {
		case "enter", "esc":
			setField(w, m.fieldCursor, m.textBuf)
			m.editingText = false
		case "backspace":
			if len(m.textBuf) > 0 {
				runes := []rune(m.textBuf)
				m.textBuf = string(runes[:len(runes)-1])
			}
		default:
			s := key.String()
			if len([]rune(s)) == 1 {
				m.textBuf += s
			}
		}
		return m, nil
	}

	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenLine
	case "j", "down", "tab":
		m.fieldCursor++
		if m.fieldCursor >= len(fields) {
			m.fieldCursor = 0
		}
	case "k", "up":
		m.fieldCursor--
		if m.fieldCursor < 0 {
			m.fieldCursor = len(fields) - 1
		}
	case "enter", " ":
		if m.fieldCursor < len(fields) {
			if fields[m.fieldCursor].isText {
				m.textBuf = getField(w, m.fieldCursor)
				m.editingText = true
			} else {
				toggleBool(w, m.fieldCursor)
			}
		}
	case "c":
		if m.fieldCursor < len(fields) &&
			(fields[m.fieldCursor].label == "fg" || fields[m.fieldCursor].label == "bg") {
			m.colorPickerTarget = fields[m.fieldCursor].label
			m.colorCursor = 0
			m.screen = screenColorPicker
		}
	case "ctrl+s":
		m = m.save()
		m.screen = screenLine
	}
	return m, nil
}

func (m model) updateTypeSelect(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenLine
	case "j", "down":
		m.typeCursor++
		if m.typeCursor >= len(allWidgetTypes) {
			m.typeCursor = 0
		}
	case "k", "up":
		m.typeCursor--
		if m.typeCursor < 0 {
			m.typeCursor = len(allWidgetTypes) - 1
		}
	case "enter":
		t := allWidgetTypes[m.typeCursor]
		id := fmt.Sprintf("%d", time.Now().UnixNano())
		w := config.WidgetItem{ID: id, Type: t}
		m.settings.Lines[m.linesCursor] = append(
			m.settings.Lines[m.linesCursor], w,
		)
		m.widgetsCursor = len(m.settings.Lines[m.linesCursor]) - 1
		m.screen = screenLine
	}
	return m, nil
}

// ---------------------------------------------------------------------------
// Global Options helpers
// ---------------------------------------------------------------------------

const globalFieldCount = 5

func isGlobalFieldText(idx int) bool {
	return idx < 3 // 0=flexMode, 1=colorLevel, 2=compactThreshold are text; 3,4 are bool
}

func getGlobalField(s *config.Settings, idx int) string {
	switch idx {
	case 0:
		return s.FlexMode
	case 1:
		return strconv.Itoa(s.ColorLevel)
	case 2:
		return strconv.Itoa(s.CompactThreshold)
	case 3:
		if s.GlobalBold {
			return "true"
		}
		return "false"
	case 4:
		if s.InheritSeparatorColors {
			return "true"
		}
		return "false"
	}
	return ""
}

func setGlobalField(s *config.Settings, idx int, val string) {
	switch idx {
	case 0:
		s.FlexMode = val
	case 1:
		if n, err := strconv.Atoi(val); err == nil {
			s.ColorLevel = n
		}
	case 2:
		if n, err := strconv.Atoi(val); err == nil {
			s.CompactThreshold = n
		}
	case 3:
		s.GlobalBold = val == "true"
	case 4:
		s.InheritSeparatorColors = val == "true"
	}
}

func toggleGlobalField(s *config.Settings, idx int) {
	switch idx {
	case 3:
		s.GlobalBold = !s.GlobalBold
	case 4:
		s.InheritSeparatorColors = !s.InheritSeparatorColors
	}
}

var globalFieldLabels = [globalFieldCount]string{
	"flexMode",
	"colorLevel",
	"compactThreshold",
	"globalBold",
	"inheritSeparatorColors",
}

func (m model) updateGlobalOptions(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editingText {
		switch key.String() {
		case "enter", "esc":
			setGlobalField(m.settings, m.fieldCursor, m.textBuf)
			m.editingText = false
		case "backspace":
			if len(m.textBuf) > 0 {
				runes := []rune(m.textBuf)
				m.textBuf = string(runes[:len(runes)-1])
			}
		default:
			s := key.String()
			if len([]rune(s)) == 1 {
				m.textBuf += s
			}
		}
		return m, nil
	}

	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenLines
		m.fieldCursor = 0
	case "j", "down", "tab":
		m.fieldCursor++
		if m.fieldCursor >= globalFieldCount {
			m.fieldCursor = 0
		}
	case "k", "up":
		m.fieldCursor--
		if m.fieldCursor < 0 {
			m.fieldCursor = globalFieldCount - 1
		}
	case "enter", " ":
		if isGlobalFieldText(m.fieldCursor) {
			m.textBuf = getGlobalField(m.settings, m.fieldCursor)
			m.editingText = true
		} else {
			toggleGlobalField(m.settings, m.fieldCursor)
		}
	case "ctrl+s":
		m = m.save()
		m.screen = screenLines
	}
	return m, nil
}

// ---------------------------------------------------------------------------
// Powerline helpers
// ---------------------------------------------------------------------------

const powerlineFieldCount = 6

func isPowerlineFieldText(idx int) bool {
	// 0=enabled(bool), 1=separator(text), 2=startCap(text), 3=endCap(text), 4=autoAlign(bool), 5=theme(text)
	return (idx >= 1 && idx <= 3) || idx == 5
}

func getPowerlineField(s *config.Settings, idx int) string {
	if s.Powerline == nil {
		switch idx {
		case 0:
			return "false"
		case 4:
			return "false"
		}
		return ""
	}
	switch idx {
	case 0:
		if s.Powerline.Enabled {
			return "true"
		}
		return "false"
	case 1:
		if len(s.Powerline.Separators) > 0 {
			return s.Powerline.Separators[0]
		}
		return ""
	case 2:
		if len(s.Powerline.StartCaps) > 0 {
			return s.Powerline.StartCaps[0]
		}
		return ""
	case 3:
		if len(s.Powerline.EndCaps) > 0 {
			return s.Powerline.EndCaps[0]
		}
		return ""
	case 4:
		if s.Powerline.AutoAlign {
			return "true"
		}
		return "false"
	case 5:
		return s.Powerline.Theme
	}
	return ""
}

func setPowerlineField(s *config.Settings, idx int, val string) {
	if s.Powerline == nil {
		s.Powerline = &config.PowerlineConfig{}
	}
	switch idx {
	case 0:
		s.Powerline.Enabled = val == "true"
	case 1:
		if len(s.Powerline.Separators) == 0 {
			s.Powerline.Separators = []string{val}
		} else {
			s.Powerline.Separators[0] = val
		}
	case 2:
		if len(s.Powerline.StartCaps) == 0 {
			s.Powerline.StartCaps = []string{val}
		} else {
			s.Powerline.StartCaps[0] = val
		}
	case 3:
		if len(s.Powerline.EndCaps) == 0 {
			s.Powerline.EndCaps = []string{val}
		} else {
			s.Powerline.EndCaps[0] = val
		}
	case 4:
		s.Powerline.AutoAlign = val == "true"
	case 5:
		s.Powerline.Theme = val
	}
}

func togglePowerlineField(s *config.Settings, idx int) {
	if s.Powerline == nil {
		s.Powerline = &config.PowerlineConfig{}
	}
	switch idx {
	case 0:
		s.Powerline.Enabled = !s.Powerline.Enabled
	case 4:
		s.Powerline.AutoAlign = !s.Powerline.AutoAlign
	}
}

var powerlineFieldLabels = [powerlineFieldCount]string{
	"enabled",
	"separator",
	"startCap",
	"endCap",
	"autoAlign",
	"theme",
}

// ---------------------------------------------------------------------------
// Install helpers
// ---------------------------------------------------------------------------

// findCCStatuslineBinary attempts to locate the ccstatusline renderer binary.
// It first checks the same directory as the running executable, then falls
// back to PATH lookup.
func findCCStatuslineBinary() string {
	// 1. Same dir as current executable — replace "ccstatusline-setup" → "ccstatusline".
	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), "ccstatusline")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	// 2. PATH lookup.
	if path, err := exec.LookPath("ccstatusline"); err == nil {
		return path
	}
	return ""
}

func (m model) updateInstall(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenLines
	case "i":
		path := findCCStatuslineBinary()
		if path == "" {
			m.statusMsg = "ccstatusline binary not found — build first with: go build -o bin/ccstatusline ./cmd/ccstatusline"
		} else if err := config.Install(path); err != nil {
			m.statusMsg = "install error: " + err.Error()
		} else {
			m.statusMsg = "installed in Claude Code settings."
		}
	case "u":
		if err := config.Uninstall(); err != nil {
			m.statusMsg = "uninstall error: " + err.Error()
		} else {
			m.statusMsg = "uninstalled from Claude Code settings."
		}
	}
	return m, nil
}

// ---------------------------------------------------------------------------
// Preview helpers
// ---------------------------------------------------------------------------

// previewInput returns a synthetic StatusInput for the live preview.
func previewInput() *input.StatusInput {
	used := 42.0
	remaining := 58.0
	cost := 0.0123
	dur := int64(150000)
	return &input.StatusInput{
		HookEventName: "Status",
		Model: &input.ModelInfo{
			ID:          "claude-opus-4-6[1m]",
			DisplayName: "Opus 4.6 (1M context)",
		},
		Version: "2.1.80",
		Cost: &input.CostInfo{
			TotalCostUSD:    cost,
			TotalDurationMs: dur,
		},
		ContextWindow: &input.ContextWindow{
			TotalInputTokens:  420000,
			TotalOutputTokens: 20000,
			ContextWindowSize: 1000000,
			UsedPercentage:    used,
			RemainingPct:      remaining,
		},
	}
}

func (m model) renderPreview() []string {
	si := previewInput()
	return render.Render(si, m.settings, 220)
}

func (m model) updatePreview(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenLines
	case "r":
		m.previewLines = m.renderPreview()
	}
	return m, nil
}

// ---------------------------------------------------------------------------
// Color picker helpers
// ---------------------------------------------------------------------------

func (m model) updateColorPicker(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenWidget
	case "j", "down":
		m.colorCursor++
		if m.colorCursor >= len(colorOptions) {
			m.colorCursor = 0
		}
	case "k", "up":
		m.colorCursor--
		if m.colorCursor < 0 {
			m.colorCursor = len(colorOptions) - 1
		}
	case "enter":
		w := m.currentWidget()
		if w != nil {
			fields := fieldsFor(w)
			// Find the field index matching colorPickerTarget.
			for i, f := range fields {
				if f.label == m.colorPickerTarget {
					setField(w, i, colorOptions[m.colorCursor].code)
					break
				}
			}
		}
		m.screen = screenWidget
	}
	return m, nil
}

func (m model) updatePowerline(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editingText {
		switch key.String() {
		case "enter", "esc":
			setPowerlineField(m.settings, m.fieldCursor, m.textBuf)
			m.editingText = false
		case "backspace":
			if len(m.textBuf) > 0 {
				runes := []rune(m.textBuf)
				m.textBuf = string(runes[:len(runes)-1])
			}
		default:
			s := key.String()
			if len([]rune(s)) == 1 {
				m.textBuf += s
			}
		}
		return m, nil
	}

	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.screen = screenLines
		m.fieldCursor = 0
	case "j", "down", "tab":
		m.fieldCursor++
		if m.fieldCursor >= powerlineFieldCount {
			m.fieldCursor = 0
		}
	case "k", "up":
		m.fieldCursor--
		if m.fieldCursor < 0 {
			m.fieldCursor = powerlineFieldCount - 1
		}
	case "enter", " ":
		if isPowerlineFieldText(m.fieldCursor) {
			m.textBuf = getPowerlineField(m.settings, m.fieldCursor)
			m.editingText = true
		} else {
			togglePowerlineField(m.settings, m.fieldCursor)
		}
	case "ctrl+s":
		m = m.save()
		m.screen = screenLines
	}
	return m, nil
}
