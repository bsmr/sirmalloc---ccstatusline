package widgets_test

import (
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func TestContextPctUsableWidget(t *testing.T) {
	w := mustGetWidget(t, "context-pct-usable")
	item := config.WidgetItem{Type: "context-pct-usable"}

	tests := []struct {
		name  string
		si    *input.StatusInput
		want  string
	}{
		{
			name: "nil StatusInput",
			si:   nil,
			want: "",
		},
		{
			name: "nil ContextWindow",
			si:   &input.StatusInput{},
			want: "",
		},
		{
			// 1M model → usable = 800k. 400k used → (800k-400k)/800k = 50%
			name: "1M model 400k used",
			si: &input.StatusInput{
				Model: &input.ModelInfo{ID: "claude-opus-4-6[1m]"},
				ContextWindow: &input.ContextWindow{
					TotalInputTokens:  400_000,
					TotalOutputTokens: 0,
				},
			},
			want: "50%",
		},
		{
			// default model → usable = 160k. 0 used → 100%
			name: "200k model 0 used",
			si: &input.StatusInput{
				Model: &input.ModelInfo{ID: "claude-sonnet"},
				ContextWindow: &input.ContextWindow{
					TotalInputTokens:  0,
					TotalOutputTokens: 0,
				},
			},
			want: "100%",
		},
		{
			// over limit: 1M model, 900k used → (800k-900k)/800k = -12.5% → clamped to 0%
			name: "1M model over usable limit",
			si: &input.StatusInput{
				Model: &input.ModelInfo{ID: "claude-opus-4-6[1m]"},
				ContextWindow: &input.ContextWindow{
					TotalInputTokens:  900_000,
					TotalOutputTokens: 0,
				},
			},
			want: "0%",
		},
		{
			// no model info → usable = 160k (default). 160k used → 0%
			name: "default model fully used",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{
					TotalInputTokens:  160_000,
					TotalOutputTokens: 0,
				},
			},
			want: "0%",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := w.Render(tc.si, item)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
