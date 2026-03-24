package render_test

import (
	"strings"
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/render"

	// Register all widgets.
	_ "go.a8l.eu/ccstatusline/internal/widgets"
)

func exampleInput() *input.StatusInput {
	cost := 0.01234
	_ = cost
	usedPct := float64(42)
	resetsAt := int64(1774020000)
	return &input.StatusInput{
		HookEventName: "Status",
		Model: &input.ModelInfo{
			ID:          "claude-opus-4-6[1m]",
			DisplayName: "Opus 4.6 (1M context)",
		},
		Cost: &input.CostInfo{
			TotalCostUSD:    0.01234,
			TotalDurationMs: 45000,
		},
		ContextWindow: &input.ContextWindow{
			UsedPercentage: 8,
			RemainingPct:   92,
		},
		RateLimits: &input.RateLimits{
			FiveHour: &input.RateLimitPeriod{
				UsedPercentage: &usedPct,
				ResetsAt:       &resetsAt,
			},
		},
	}
}

func TestRenderDefaultSettings(t *testing.T) {
	si := exampleInput()
	cfg := config.DefaultSettings()

	lines := render.Render(si, cfg, 80)

	if len(lines) == 0 {
		t.Fatal("Render returned no lines")
	}
	line := lines[0]
	if line == "" {
		t.Fatal("Render returned empty first line")
	}

	// Must contain model name, context %, duration, cost.
	checks := []string{
		"Opus 4.6 (1M context)",
		"8%",
		"0m",
		"$0.0123",
	}
	for _, want := range checks {
		if !strings.Contains(line, want) {
			t.Errorf("output %q does not contain %q", line, want)
		}
	}
}

func TestRenderNilFields(t *testing.T) {
	// All pointer fields nil — must not panic.
	si := &input.StatusInput{}
	cfg := config.DefaultSettings()

	lines := render.Render(si, cfg, 80)
	if lines == nil {
		t.Fatal("Render returned nil")
	}
}
