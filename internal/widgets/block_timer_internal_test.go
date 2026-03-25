// White-box tests for unexported block-timer helpers.
package widgets

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input" //nolint:depguard
)

// ---------------------------------------------------------------------------
// parseBlockStart
// ---------------------------------------------------------------------------

func TestParseBlockStart_ValidBoundary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "transcript.jsonl")

	// hour=10, 10%5==0 → valid block boundary
	blockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	line := fmt.Sprintf(`{"timestamp":%q}`, blockTime.Format(time.RFC3339))
	if err := os.WriteFile(path, []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := parseBlockStart(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(blockTime) {
		t.Errorf("expected %v, got %v", blockTime, got)
	}
}

func TestParseBlockStart_NoBlockBoundary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "transcript.jsonl")

	// hour=11, 11%5!=0 → no block start
	ts := time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC)
	line := fmt.Sprintf(`{"timestamp":%q}`, ts.Format(time.RFC3339))
	if err := os.WriteFile(path, []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := parseBlockStart(path); err == nil {
		t.Error("expected error when no block boundary found")
	}
}

func TestParseBlockStart_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseBlockStart(path); err == nil {
		t.Error("expected error for empty file")
	}
}

func TestParseBlockStart_MalformedLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "transcript.jsonl")
	content := "not json\n{\"no_timestamp\":true}\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseBlockStart(path); err == nil {
		t.Error("expected error for file with no parseable timestamps")
	}
}

func TestParseBlockStart_MixedLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "transcript.jsonl")

	// First entry at 11:30 (no boundary), second at 15:00 (15%5==0).
	t1 := time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)
	content := fmt.Sprintf("%s\n%s\n",
		fmt.Sprintf(`{"timestamp":%q}`, t1.Format(time.RFC3339)),
		fmt.Sprintf(`{"timestamp":%q}`, t2.Format(time.RFC3339)),
	)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := parseBlockStart(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(t2) {
		t.Errorf("expected last boundary %v, got %v", t2, got)
	}
}

func TestParseBlockStart_FileNotFound(t *testing.T) {
	if _, err := parseBlockStart("/nonexistent/path/transcript.jsonl"); err == nil {
		t.Error("expected error for missing file")
	}
}

// ---------------------------------------------------------------------------
// renderBar
// ---------------------------------------------------------------------------

func TestRenderBar(t *testing.T) {
	tests := []struct {
		elapsed  time.Duration
		total    time.Duration
		width    int
		wantFill int // number of filled runes
	}{
		{0, 5 * time.Hour, 10, 0},
		{5 * time.Hour, 5 * time.Hour, 10, 10},
		{2*time.Hour + 30*time.Minute, 5 * time.Hour, 10, 5},
		{0, 5 * time.Hour, 5, 0},
		{5 * time.Hour, 5 * time.Hour, 5, 5},
		{6 * time.Hour, 5 * time.Hour, 10, 10}, // over 100% → clamped
		{-1, 5 * time.Hour, 10, 0},             // negative → clamped
	}
	for _, tc := range tests {
		got := renderBar(tc.elapsed, tc.total, tc.width)
		runes := []rune(got)
		if len(runes) != tc.width {
			t.Errorf("renderBar(%v, %v, %d): len=%d, want %d (got %q)",
				tc.elapsed, tc.total, tc.width, len(runes), tc.width, got)
			continue
		}
		filled := 0
		for _, r := range runes {
			if r == '█' {
				filled++
			}
		}
		if filled != tc.wantFill {
			t.Errorf("renderBar(%v, %v, %d): filled=%d, want %d (got %q)",
				tc.elapsed, tc.total, tc.width, filled, tc.wantFill, got)
		}
	}
}

// ---------------------------------------------------------------------------
// getBlockStart — path validation
// ---------------------------------------------------------------------------

func TestGetBlockStart_NilInput(t *testing.T) {
	if _, err := getBlockStart(nil); err == nil {
		t.Error("expected error for nil input")
	}
}

func TestGetBlockStart_EmptyTranscriptPath(t *testing.T) {
	if _, err := getBlockStart(&input.StatusInput{}); err == nil {
		t.Error("expected error for empty transcript path")
	}
}

func TestGetBlockStart_RelativePath(t *testing.T) {
	si := &input.StatusInput{TranscriptPath: "relative/path.jsonl"}
	if _, err := getBlockStart(si); err == nil {
		t.Error("expected error for relative path")
	}
}

func TestGetBlockStart_PathOutsideHome(t *testing.T) {
	si := &input.StatusInput{TranscriptPath: "/tmp/transcript.jsonl"}
	if _, err := getBlockStart(si); err == nil {
		t.Error("expected error for path outside home dir")
	}
}

func TestGetBlockStart_InvalidSessionID(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	dir, err := os.MkdirTemp(home, ".ccstatusline-test-*")
	if err != nil {
		t.Skip("cannot create temp dir in home:", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	transcriptPath := filepath.Join(dir, "transcript.jsonl")
	if err := os.WriteFile(transcriptPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	si := &input.StatusInput{
		TranscriptPath: transcriptPath,
		SessionID:      "../../etc/passwd", // invalid characters
	}
	if _, err := getBlockStart(si); err == nil {
		t.Error("expected error for invalid session ID")
	}
}

// ---------------------------------------------------------------------------
// blockTimerWidget.Render — integration with real transcript
// ---------------------------------------------------------------------------

// transcriptUnderHome creates a valid transcript JSONL under the user's home
// and returns the path and a cleanup function. Skips the test if home dir
// is not writable.
func transcriptUnderHome(t *testing.T, blockHour int) string {
	t.Helper()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	dir, err := os.MkdirTemp(home, ".ccstatusline-test-*")
	if err != nil {
		t.Skip("cannot create temp dir in home:", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	blockTime := time.Date(2024, 1, 15, blockHour, 0, 0, 0, time.UTC)
	path := filepath.Join(dir, "transcript.jsonl")
	line := fmt.Sprintf(`{"timestamp":%q}`, blockTime.Format(time.RFC3339))
	if err := os.WriteFile(path, []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestBlockTimerWidget_Render_DefaultMode(t *testing.T) {
	path := transcriptUnderHome(t, 10) // 10%5==0
	w := blockTimerWidget{}
	si := &input.StatusInput{TranscriptPath: path}
	got := w.Render(si, config.WidgetItem{Type: "block-timer"})
	// Result is a duration string — just verify it is non-empty and has no panic.
	if got == "" {
		t.Error("expected non-empty duration string")
	}
}

func TestBlockTimerWidget_Render_BarMode(t *testing.T) {
	path := transcriptUnderHome(t, 10)
	w := blockTimerWidget{}
	si := &input.StatusInput{TranscriptPath: path}

	got := w.Render(si, config.WidgetItem{Type: "block-timer", BarMode: "bar"})
	if runes := []rune(got); len(runes) != 10 {
		t.Errorf("bar mode: expected 10-char bar, got %q (len=%d)", got, len(runes))
	}
}

func TestBlockTimerWidget_Render_BarShortMode(t *testing.T) {
	path := transcriptUnderHome(t, 10)
	w := blockTimerWidget{}
	si := &input.StatusInput{TranscriptPath: path}

	got := w.Render(si, config.WidgetItem{Type: "block-timer", BarMode: "bar-short"})
	if runes := []rune(got); len(runes) != 5 {
		t.Errorf("bar-short mode: expected 5-char bar, got %q (len=%d)", got, len(runes))
	}
}

func TestBlockTimerWidget_Render_CacheHit(t *testing.T) {
	path := transcriptUnderHome(t, 10)
	w := blockTimerWidget{}
	si := &input.StatusInput{TranscriptPath: path}
	item := config.WidgetItem{Type: "block-timer"}

	// First call writes cache; second call should use it.
	r1 := w.Render(si, item)
	r2 := w.Render(si, item)
	if r1 == "" || r2 == "" {
		t.Error("expected non-empty results from both calls")
	}
}
