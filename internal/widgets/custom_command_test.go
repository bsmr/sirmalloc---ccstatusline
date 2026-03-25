package widgets_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/widgets"
)

// makeTempScript creates a temporary executable shell script containing body.
func makeTempScript(t *testing.T, body string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-cmd-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("#!/bin/sh\n" + body + "\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Chmod(0o755); err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return name
}

func TestCustomCommandWidget_EmptyCommand(t *testing.T) {
	w, _ := widgets.Get("custom-command")
	got := w.Render(nil, config.WidgetItem{Type: "custom-command"})
	if got != "" {
		t.Errorf("expected empty for empty command, got %q", got)
	}
}

func TestCustomCommandWidget_CommandNotFound(t *testing.T) {
	w, _ := widgets.Get("custom-command")
	item := config.WidgetItem{Type: "custom-command", Command: "/nonexistent/command-xyz"}
	got := w.Render(nil, item)
	if got != "" {
		t.Errorf("expected empty for missing command, got %q", got)
	}
}

func TestCustomCommandWidget_OutputTrimmed(t *testing.T) {
	script := makeTempScript(t, "printf 'hello world'")
	w, _ := widgets.Get("custom-command")
	item := config.WidgetItem{Type: "custom-command", Command: script}
	got := w.Render(nil, item)
	if got != "hello world" {
		t.Errorf(`expected "hello world", got %q`, got)
	}
}

func TestCustomCommandWidget_TrailingNewlineTrimmed(t *testing.T) {
	script := makeTempScript(t, "printf 'trimmed\\n'")
	w, _ := widgets.Get("custom-command")
	item := config.WidgetItem{Type: "custom-command", Command: script}
	got := w.Render(nil, item)
	if got != "trimmed" {
		t.Errorf(`expected "trimmed", got %q`, got)
	}
}

func TestCustomCommandWidget_NonZeroExitReturnsEmpty(t *testing.T) {
	script := makeTempScript(t, "exit 1")
	w, _ := widgets.Get("custom-command")
	item := config.WidgetItem{Type: "custom-command", Command: script}
	got := w.Render(nil, item)
	if got != "" {
		t.Errorf("expected empty for non-zero exit, got %q", got)
	}
}

func TestCustomCommandWidget_Timeout(t *testing.T) {
	// Use "yes" (outputs "y\n" in a loop, single process, no fork) so the
	// context timeout kills the process directly and Output() returns promptly.
	yesPath, err := exec.LookPath("yes")
	if err != nil {
		t.Skip("yes not available:", err)
	}
	w, _ := widgets.Get("custom-command")
	item := config.WidgetItem{
		Type:    "custom-command",
		Command: yesPath,
		Timeout: 50, // 50ms
	}
	start := time.Now()
	got := w.Render(nil, item)
	elapsed := time.Since(start)

	if got != "" {
		t.Errorf("expected empty on timeout, got %q", got)
	}
	if elapsed > 2*time.Second {
		t.Errorf("timeout did not activate: elapsed %v", elapsed)
	}
}

func TestCustomCommandWidget_DefaultTimeout(t *testing.T) {
	// Timeout=0 → default 2s is used; a fast command should still succeed.
	script := makeTempScript(t, "printf 'fast'")
	w, _ := widgets.Get("custom-command")
	item := config.WidgetItem{Type: "custom-command", Command: script, Timeout: 0}
	got := w.Render(nil, item)
	if got != "fast" {
		t.Errorf(`expected "fast", got %q`, got)
	}
}
