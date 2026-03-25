package tui

import (
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.a8l.eu/ccstatusline/internal/config"
)

type screenKind int

const (
	screenLines      screenKind = iota
	screenLine
	screenWidget
	screenTypeSelect
)

// allWidgetTypes lists all widget types available for addition.
var allWidgetTypes = []string{
	"model",
	"context-percentage",
	"context-pct-usable",
	"context-length",
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
