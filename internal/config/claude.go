package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// claudeStatusLineEntry is the JSON structure written under the "statusLine" key.
type claudeStatusLineEntry struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Padding int    `json:"padding"`
}

// ClaudeConfigDir returns the Claude config directory.
// Respects $CLAUDE_CONFIG_DIR (R-6: must be absolute path).
func ClaudeConfigDir() (string, error) {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		clean := filepath.Clean(dir)
		if !filepath.IsAbs(clean) {
			return "", fmt.Errorf("CLAUDE_CONFIG_DIR must be an absolute path, got %q", dir)
		}
		return clean, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home dir: %w", err)
	}
	return filepath.Join(home, ".claude"), nil
}

// ClaudeSettingsPath returns the path to ~/.claude/settings.json.
func ClaudeSettingsPath() (string, error) {
	dir, err := ClaudeConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

// IsInstalled reports whether ccstatusline is installed in Claude Code settings.
// Returns (installed, currentCommand, error).
func IsInstalled() (bool, string, error) {
	path, err := ClaudeSettingsPath()
	if err != nil {
		return false, "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("reading Claude settings %q: %w", path, err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return false, "", fmt.Errorf("parsing Claude settings %q: %w", path, err)
	}

	slRaw, ok := raw["statusLine"]
	if !ok {
		return false, "", nil
	}

	var entry claudeStatusLineEntry
	if err := json.Unmarshal(slRaw, &entry); err != nil {
		return false, "", nil
	}

	if strings.Contains(entry.Command, "ccstatusline") {
		return true, entry.Command, nil
	}
	return false, entry.Command, nil
}

// Install writes the statusLine entry into Claude Code settings.
// commandPath is the absolute path to the ccstatusline binary.
// Reads the existing settings.json (if any), merges the statusLine entry,
// and writes back atomically (R-4).
func Install(commandPath string) error {
	path, err := ClaudeSettingsPath()
	if err != nil {
		return err
	}

	// Read existing settings (or start fresh).
	raw := make(map[string]json.RawMessage)
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading Claude settings %q: %w", path, err)
	}
	if err == nil {
		if err := json.Unmarshal(data, &raw); err != nil {
			return fmt.Errorf("parsing Claude settings %q: %w", path, err)
		}
	}

	// Build and encode the statusLine entry.
	entry := claudeStatusLineEntry{
		Type:    "command",
		Command: commandPath,
		Padding: 0,
	}
	entryBytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("encoding statusLine entry: %w", err)
	}
	raw["statusLine"] = json.RawMessage(entryBytes)

	// Marshal the merged map.
	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling Claude settings: %w", err)
	}

	return writeClaudeSettingsAtomic(path, out)
}

// Uninstall removes the statusLine entry from Claude Code settings.
func Uninstall() error {
	path, err := ClaudeSettingsPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // nothing to uninstall
		}
		return fmt.Errorf("reading Claude settings %q: %w", path, err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parsing Claude settings %q: %w", path, err)
	}

	if _, ok := raw["statusLine"]; !ok {
		return nil // already absent
	}
	delete(raw, "statusLine")

	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling Claude settings: %w", err)
	}

	return writeClaudeSettingsAtomic(path, out)
}

// writeClaudeSettingsAtomic writes data to path atomically via temp-file + rename (R-4).
func writeClaudeSettingsAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating Claude config dir: %w", err)
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
