package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	Version                int              `json:"version"`
	Lines                  [][]WidgetItem   `json:"lines"`
	FlexMode               string           `json:"flexMode"`   // "full", "full-minus-40", "full-until-compact"
	CompactThreshold       int              `json:"compactThreshold"`
	ColorLevel             int              `json:"colorLevel"` // 0=basic, 1=256, 2=truecolor
	InheritSeparatorColors bool             `json:"inheritSeparatorColors"`
	GlobalBold             bool             `json:"globalBold"`
	Powerline              *PowerlineConfig `json:"powerline,omitempty"`
}

type WidgetItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`       // kebab-case: "model", "context-percentage", "session-clock", etc.
	Color     string `json:"color,omitempty"`
	FG        string `json:"fg,omitempty"`
	BG        string `json:"bg,omitempty"`
	Bold      bool   `json:"bold,omitempty"`
	RawValue  bool   `json:"rawValue,omitempty"`
	Merge     bool   `json:"merge,omitempty"`
	Text      string `json:"text,omitempty"`
	Command   string `json:"command,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
	Segments  int    `json:"segments,omitempty"`
	FishStyle bool   `json:"fishStyle,omitempty"`
	HideNoGit bool   `json:"hideNoGit,omitempty"`
	Remaining       bool              `json:"remaining,omitempty"`
	BarMode         string            `json:"barMode,omitempty"`
	SepChar         string            `json:"sepChar,omitempty"`
	BackgroundColor string            `json:"backgroundColor,omitempty"` // e.g. "bgCyan", "bgWhite"
	Metadata        map[string]string `json:"metadata,omitempty"`         // widget-specific options
}

type PowerlineConfig struct {
	Enabled                   bool     `json:"enabled"`
	Separators                []string `json:"separators"`
	SeparatorInvertBackground []bool   `json:"separatorInvertBackground"`
	StartCaps                 []string `json:"startCaps"`
	EndCaps                   []string `json:"endCaps"`
	Theme                     string   `json:"theme,omitempty"`
	AutoAlign                 bool     `json:"autoAlign"`
}

func DefaultSettings() *Settings {
	return &Settings{
		Version:    3,
		FlexMode:   "full-minus-40",
		ColorLevel: 2,
		Lines: [][]WidgetItem{
			{
				{ID: "1", Type: "model"},
				{ID: "2", Type: "separator"},
				{ID: "3", Type: "context-percentage"},
				{ID: "4", Type: "separator"},
				{ID: "5", Type: "session-clock"},
				{ID: "6", Type: "separator"},
				{ID: "7", Type: "session-cost"},
			},
		},
	}
}

func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolving config dir: %w", err)
	}
	return filepath.Join(dir, "ccstatusline", "settings.json"), nil
}

func Load() (*Settings, error) {
	path, err := ConfigPath()
	if err != nil {
		return DefaultSettings(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultSettings(), nil
		}
		return nil, fmt.Errorf("reading settings %q: %w", path, err)
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing settings %q: %w", path, err)
	}
	return &s, nil
}

// Save writes settings atomically via temp-file + rename (Security R-4).
func Save(path string, s *Settings) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling settings: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".settings-*.json.tmp")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // clean up on failure

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("renaming temp file: %w", err)
	}
	return nil
}
