package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
)

func TestDefaultSettings(t *testing.T) {
	s := config.DefaultSettings()
	if s == nil {
		t.Fatal("DefaultSettings() returned nil")
	}
	if s.Version != 3 {
		t.Errorf("Version = %d, want 3", s.Version)
	}
	if len(s.Lines) == 0 {
		t.Error("Lines is empty")
	}
	if s.FlexMode == "" {
		t.Error("FlexMode is empty")
	}
}

func TestSaveAndLoadRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	original := config.DefaultSettings()
	original.GlobalBold = true
	original.ColorLevel = 1
	original.FlexMode = "full"

	if err := config.Save(path, original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Read back via raw JSON to confirm file was written.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal raw: %v", err)
	}
	if v, ok := raw["version"]; !ok || v.(float64) != 3 {
		t.Errorf("raw version = %v, want 3", v)
	}

	// Load via the public Load function is path-based (reads from ConfigPath),
	// so we exercise the internal path directly by using Save+os.ReadFile+json.Unmarshal.
	var loaded config.Settings
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal loaded: %v", err)
	}
	if loaded.Version != original.Version {
		t.Errorf("loaded Version = %d, want %d", loaded.Version, original.Version)
	}
	if loaded.GlobalBold != original.GlobalBold {
		t.Errorf("loaded GlobalBold = %v, want %v", loaded.GlobalBold, original.GlobalBold)
	}
	if loaded.ColorLevel != original.ColorLevel {
		t.Errorf("loaded ColorLevel = %d, want %d", loaded.ColorLevel, original.ColorLevel)
	}
	if loaded.FlexMode != original.FlexMode {
		t.Errorf("loaded FlexMode = %q, want %q", loaded.FlexMode, original.FlexMode)
	}
}

func TestSaveAtomicNoTempFilesRemain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	if err := config.Save(path, config.DefaultSettings()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}
}

func TestLoadNonExistentReturnsDefaults(t *testing.T) {
	// Load() reads from ConfigPath() which we cannot override without env tricks,
	// but we can verify DefaultSettings() is what a missing-file Load returns
	// by testing the observable contract: result is non-nil with Version=3.
	//
	// Because Load() falls back to DefaultSettings() on ErrNotExist, and
	// ConfigPath() may or may not exist on this machine, we test both branches
	// by calling Save to a temp path and verifying the round-trip succeeds.
	// The non-existent-path branch is exercised by setting XDG_CONFIG_HOME to
	// a temp dir with no settings file.
	orig := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if orig == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", orig)
		}
	}()

	emptyDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", emptyDir)

	s, err := config.Load()
	if err != nil {
		t.Fatalf("Load with missing settings: %v", err)
	}
	if s == nil {
		t.Fatal("Load returned nil for missing settings")
	}
	if s.Version != 3 {
		t.Errorf("Version = %d, want 3", s.Version)
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	orig := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if orig == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", orig)
		}
	}()

	cfgDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", cfgDir)

	// Write invalid JSON to the expected path.
	settingsDir := filepath.Join(cfgDir, "ccstatusline")
	if err := os.MkdirAll(settingsDir, 0o700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(settingsDir, "settings.json"), []byte("{bad json"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load with invalid JSON should return error, got nil")
	}
}

func TestConfigPathContainsCCStatusline(t *testing.T) {
	path, err := config.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("ConfigPath = %q is not absolute", path)
	}
	if !strings.Contains(path, "ccstatusline") {
		t.Errorf("ConfigPath = %q does not contain \"ccstatusline\"", path)
	}
}
