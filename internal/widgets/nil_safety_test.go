package widgets_test

import (
	"fmt"
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/widgets"
)

// allRegisteredTypes lists every widget type registered by this package.
// Update this list when new widget types are added.
var allRegisteredTypes = []string{
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
}

func TestNilSafetyAllWidgets(t *testing.T) {
	for _, typ := range allRegisteredTypes {
		typ := typ
		t.Run(typ, func(t *testing.T) {
			w, ok := widgets.Get(typ)
			if !ok {
				t.Skipf("widget type %q not registered — skipping nil-safety check", typ)
			}

			var result string
			var panicked any

			func() {
				defer func() {
					panicked = recover()
				}()
				result = w.Render(nil, config.WidgetItem{Type: typ})
			}()

			if panicked != nil {
				t.Errorf("widget %q panicked on nil StatusInput: %v", typ, panicked)
			}
			// Result must be a string (empty is acceptable).
			_ = fmt.Sprintf("%s", result) // compile-time type assertion
		})
	}
}
