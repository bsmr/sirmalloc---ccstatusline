package widgets_test

import (
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/widgets"
)

// dummyWidget is a minimal Widget implementation used in registration tests.
type dummyWidget struct{}

func (dummyWidget) Render(_ *input.StatusInput, _ config.WidgetItem) string { return "" }

func TestRegisterAndGet(t *testing.T) {
	const typ = "test-widget-roundtrip"
	widgets.Register(typ, dummyWidget{})

	w, ok := widgets.Get(typ)
	if !ok {
		t.Fatalf("Get(%q) = false, want true", typ)
	}
	if w == nil {
		t.Fatalf("Get(%q) returned nil widget", typ)
	}
}

func TestGetUnknownType(t *testing.T) {
	w, ok := widgets.Get("__nonexistent__")
	if ok {
		t.Error("Get with unknown type returned ok=true, want false")
	}
	if w != nil {
		t.Error("Get with unknown type returned non-nil widget")
	}
}

func TestExpectedWidgetTypesAreRegistered(t *testing.T) {
	expected := []string{
		"model",
		"context-percentage",
		"session-clock",
		"session-cost",
		"separator",
		"flex-separator",
		"git-branch",
		"cwd",
		"tokens-input",
		"tokens-output",
		"tokens-cached",
		"tokens-total",
		"context-length",
		"context-pct-usable",
		"custom-command",
		"custom-text",
		"block-timer",
		"vim-mode",
		"rate-limit-five-hour",
		"rate-limit-seven-day",
		"terminal-width",
		"output-style",
		"version",
		"git-worktree",
		"git-changes",
		"context-bar",
		"reset-timer",
		"session-usage",
	}
	for _, typ := range expected {
		t.Run(typ, func(t *testing.T) {
			_, ok := widgets.Get(typ)
			if !ok {
				t.Errorf("widget type %q is not registered", typ)
			}
		})
	}
}
