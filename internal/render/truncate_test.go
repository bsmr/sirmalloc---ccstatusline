package render_test

import (
	"strings"
	"testing"

	"go.a8l.eu/ccstatusline/internal/render"
)

func TestVisibleLen(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", "", 0},
		{"plain", "hello", 5},
		{"basic color", "\033[31mred\033[0m", 3},
		{"truecolor", "\033[38;2;255;0;0mred\033[0m", 3},
		{"osc8 hyperlink", "a\033]8;;http://x\007link\033]8;;\007b", 6},
		{"only escape", "\033[0m", 0},
		{"mixed", "ab\033[32mcd\033[0mef", 6},
		{"unicode multibyte", "héllo", 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := render.VisibleLen(tc.input)
			if got != tc.want {
				t.Errorf("VisibleLen(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	t.Run("no truncation needed", func(t *testing.T) {
		got := render.Truncate("hello", 10)
		if got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("exact fit", func(t *testing.T) {
		got := render.Truncate("hello", 5)
		if got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("truncate plain string", func(t *testing.T) {
		// maxVisible=5 → limit=2 visible chars + "..."
		got := render.Truncate("hello world", 5)
		if render.VisibleLen(got) > 5 {
			t.Errorf("VisibleLen(%q) = %d, want <= 5", got, render.VisibleLen(got))
		}
		if !strings.HasSuffix(got, "...") {
			t.Errorf("got %q, want suffix \"...\"", got)
		}
		// visible prefix must be "he"
		stripped := strings.ReplaceAll(got, "\033[0m", "")
		stripped = strings.TrimSuffix(stripped, "...")
		if stripped != "he" {
			t.Errorf("visible prefix = %q, want %q", stripped, "he")
		}
	})

	t.Run("truncate ansi string", func(t *testing.T) {
		// "\033[31mhello world\033[0m" with maxVisible=8 → "hello" visible + Reset + "..."
		input := "\033[31mhello world\033[0m"
		got := render.Truncate(input, 8)

		if render.VisibleLen(got) > 8 {
			t.Errorf("VisibleLen(%q) = %d, want <= 8", got, render.VisibleLen(got))
		}
		if !strings.HasPrefix(got, "\033[31m") {
			t.Errorf("got %q, want prefix \"\\033[31m\"", got)
		}
		if !strings.HasSuffix(got, "...") {
			t.Errorf("got %q, want suffix \"...\"", got)
		}
		if !strings.Contains(got, "\033[0m") {
			t.Errorf("got %q, missing Reset before ellipsis", got)
		}
		// 8 - 3 = 5 visible chars → "hello"
		if render.VisibleLen(got) != 8 {
			// visible = 5 (hello) + 3 (...) = 8
			t.Errorf("VisibleLen = %d, want 8", render.VisibleLen(got))
		}
	})

	t.Run("maxVisible zero", func(t *testing.T) {
		got := render.Truncate("hello", 0)
		if !strings.HasSuffix(got, "...") {
			t.Errorf("got %q, want suffix \"...\"", got)
		}
	})

	t.Run("maxVisible less than ellipsis", func(t *testing.T) {
		// maxVisible=2 → limit=0 → no visible chars, just Reset+"..."
		got := render.Truncate("hello", 2)
		if !strings.HasSuffix(got, "...") {
			t.Errorf("got %q, want suffix \"...\"", got)
		}
	})
}
