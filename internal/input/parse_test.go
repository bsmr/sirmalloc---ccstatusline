package input_test

import (
	"os"
	"testing"

	"go.a8l.eu/ccstatusline/internal/input"
)

func TestDecodeExampleInput(t *testing.T) {
	f, err := os.Open("../../testdata/example_input.json")
	if err != nil {
		t.Fatalf("open testdata: %v", err)
	}
	defer f.Close()

	si, err := input.Decode(f)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if si.HookEventName != "Status" {
		t.Errorf("HookEventName = %q, want %q", si.HookEventName, "Status")
	}
	if si.SessionID != "abc123..." {
		t.Errorf("SessionID = %q, want %q", si.SessionID, "abc123...")
	}

	if si.Model == nil {
		t.Fatal("Model is nil")
	}
	if si.Model.ID != "claude-opus-4-6[1m]" {
		t.Errorf("Model.ID = %q, want %q", si.Model.ID, "claude-opus-4-6[1m]")
	}
	if si.Model.DisplayName != "Opus 4.6 (1M context)" {
		t.Errorf("Model.DisplayName = %q, want %q", si.Model.DisplayName, "Opus 4.6 (1M context)")
	}

	if si.Cost == nil {
		t.Fatal("Cost is nil")
	}
	if si.Cost.TotalCostUSD != 0.01234 {
		t.Errorf("Cost.TotalCostUSD = %v, want 0.01234", si.Cost.TotalCostUSD)
	}
	if si.Cost.TotalDurationMs != 45000 {
		t.Errorf("Cost.TotalDurationMs = %v, want 45000", si.Cost.TotalDurationMs)
	}

	if si.ContextWindow == nil {
		t.Fatal("ContextWindow is nil")
	}
	if si.ContextWindow.UsedPercentage != 8 {
		t.Errorf("ContextWindow.UsedPercentage = %v, want 8", si.ContextWindow.UsedPercentage)
	}
	if si.ContextWindow.RemainingPct != 92 {
		t.Errorf("ContextWindow.RemainingPct = %v, want 92", si.ContextWindow.RemainingPct)
	}

	if si.RateLimits == nil {
		t.Fatal("RateLimits is nil")
	}
	if si.RateLimits.FiveHour == nil {
		t.Fatal("RateLimits.FiveHour is nil")
	}
	if si.RateLimits.FiveHour.UsedPercentage == nil || *si.RateLimits.FiveHour.UsedPercentage != 42 {
		t.Errorf("RateLimits.FiveHour.UsedPercentage = %v, want 42", si.RateLimits.FiveHour.UsedPercentage)
	}

	if si.Vim == nil {
		t.Fatal("Vim is nil")
	}
	if si.Vim.Mode != "NORMAL" {
		t.Errorf("Vim.Mode = %q, want %q", si.Vim.Mode, "NORMAL")
	}
}
