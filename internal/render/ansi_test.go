package render_test

import (
	"testing"

	"go.a8l.eu/ccstatusline/internal/render"
)

func TestFG(t *testing.T) {
	tests := []struct {
		color string
		want  string
	}{
		{"", ""},
		{"red", "\033[31m"},
		{"green", "\033[32m"},
		{"yellow", "\033[33m"},
		{"blue", "\033[34m"},
		{"magenta", "\033[35m"},
		{"cyan", "\033[36m"},
		{"white", "\033[37m"},
		{"black", "\033[30m"},
		// 256-color passthrough
		{"38;5;200", "\033[38;5;200m"},
		// truecolor passthrough
		{"38;2;255;128;0", "\033[38;2;255;128;0m"},
	}
	for _, tc := range tests {
		got := render.FG(tc.color)
		if got != tc.want {
			t.Errorf("FG(%q) = %q, want %q", tc.color, got, tc.want)
		}
	}
}

func TestBG(t *testing.T) {
	tests := []struct {
		color string
		want  string
	}{
		{"", ""},
		{"black", "\033[40m"},
		{"red", "\033[41m"},
		{"green", "\033[42m"},
		{"yellow", "\033[43m"},
		{"blue", "\033[44m"},
		{"magenta", "\033[45m"},
		{"cyan", "\033[46m"},
		{"white", "\033[47m"},
		// 256-color passthrough
		{"48;5;100", "\033[48;5;100m"},
	}
	for _, tc := range tests {
		got := render.BG(tc.color)
		if got != tc.want {
			t.Errorf("BG(%q) = %q, want %q", tc.color, got, tc.want)
		}
	}
}

func TestBold(t *testing.T) {
	got := render.Bold()
	if got != "\033[1m" {
		t.Errorf("Bold() = %q, want %q", got, "\033[1m")
	}
}

func TestReset(t *testing.T) {
	got := render.Reset()
	if got != "\033[0m" {
		t.Errorf("Reset() = %q, want %q", got, "\033[0m")
	}
}
