package widgets_test

import (
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/widgets"
)

func mustGetWidget(t *testing.T, typ string) widgets.Widget {
	t.Helper()
	w, ok := widgets.Get(typ)
	if !ok {
		t.Fatalf("widget type %q not registered", typ)
	}
	return w
}

func TestTokensInputWidget(t *testing.T) {
	w := mustGetWidget(t, "tokens-input")

	tests := []struct {
		name string
		si   *input.StatusInput
		want string
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
			name: "1500 tokens",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{TotalInputTokens: 1500},
			},
			want: "1.5k",
		},
		{
			name: "zero tokens",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{TotalInputTokens: 0},
			},
			want: "0",
		},
		{
			name: "1000000 tokens",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{TotalInputTokens: 1_000_000},
			},
			want: "1.0M",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := w.Render(tc.si, config.WidgetItem{Type: "tokens-input"})
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTokensTotalWidget(t *testing.T) {
	w := mustGetWidget(t, "tokens-total")

	tests := []struct {
		name string
		si   *input.StatusInput
		want string
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
			name: "in=1000 out=500",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{
					TotalInputTokens:  1000,
					TotalOutputTokens: 500,
				},
			},
			want: "1.5k",
		},
		{
			name: "both zero",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{},
			},
			want: "0",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := w.Render(tc.si, config.WidgetItem{Type: "tokens-total"})
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTokensCachedWidget(t *testing.T) {
	w := mustGetWidget(t, "tokens-cached")

	tests := []struct {
		name string
		si   *input.StatusInput
		want string
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
			name: "nil CurrentUsage",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{},
			},
			want: "",
		},
		{
			name: "2000 cache-read tokens",
			si: &input.StatusInput{
				ContextWindow: &input.ContextWindow{
					CurrentUsage: &input.CurrentUsage{
						CacheReadInputTokens: 2000,
					},
				},
			},
			want: "2.0k",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := w.Render(tc.si, config.WidgetItem{Type: "tokens-cached"})
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
