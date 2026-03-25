package widgets

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/util"
)

func init() {
	Register("block-timer", blockTimerWidget{})
}

type blockTimerWidget struct{}

// sessionIDRe validates session IDs used as filename components (Security R-3).
var sessionIDRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type blockCache struct {
	BlockStart     string `json:"block_start"`
	TranscriptPath string `json:"transcript_path"`
}

func blockCachePath(configDir string) (string, error) {
	h := sha256.Sum256([]byte(configDir))
	hash := hex.EncodeToString(h[:])[:8]
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "ccstatusline", fmt.Sprintf("block-cache-%s.json", hash)), nil
}

func loadBlockCache(cachePath string) *blockCache {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil
	}
	var c blockCache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil
	}
	return &c
}

func saveBlockCache(cachePath string, c *blockCache) {
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return
	}
	tmp, err := os.CreateTemp(dir, ".block-cache-*.json.tmp")
	if err != nil {
		return
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return
	}
	if err := tmp.Close(); err != nil {
		return
	}
	_ = os.Rename(tmpName, cachePath)
}

func parseBlockStart(transcriptPath string) (time.Time, error) {
	f, err := os.Open(transcriptPath)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	now := time.Now()
	var blockStart time.Time

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry struct {
			Timestamp string `json:"timestamp"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if entry.Timestamp == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			// try without nanoseconds
			t, err = time.Parse("2006-01-02T15:04:05Z", entry.Timestamp)
			if err != nil {
				continue
			}
		}
		if t.After(now) {
			continue
		}
		if t.Hour()%5 == 0 {
			blockStart = t
		}
	}
	if err := scanner.Err(); err != nil {
		return time.Time{}, err
	}
	if blockStart.IsZero() {
		return time.Time{}, fmt.Errorf("no block start found")
	}
	return blockStart, nil
}

func (blockTimerWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil {
		return ""
	}

	transcriptPath := si.TranscriptPath
	if transcriptPath == "" {
		return ""
	}

	// Security R-2: validate transcript path is absolute and under home dir.
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	clean := filepath.Clean(transcriptPath)
	if !filepath.IsAbs(clean) {
		return ""
	}
	if !strings.HasPrefix(clean, home+string(filepath.Separator)) {
		return ""
	}

	// Security R-3: validate session_id before using as filename part.
	if si.SessionID != "" && !sessionIDRe.MatchString(si.SessionID) {
		return ""
	}

	configDir := filepath.Dir(transcriptPath)
	cachePath, err := blockCachePath(configDir)
	if err != nil {
		return ""
	}

	var blockStart time.Time

	// Try cache first.
	if cached := loadBlockCache(cachePath); cached != nil && cached.TranscriptPath == transcriptPath {
		t, err := time.Parse(time.RFC3339, cached.BlockStart)
		if err == nil {
			blockStart = t
		}
	}

	// Parse transcript if cache miss.
	if blockStart.IsZero() {
		t, err := parseBlockStart(clean)
		if err != nil {
			return ""
		}
		blockStart = t
		saveBlockCache(cachePath, &blockCache{
			BlockStart:     blockStart.Format(time.RFC3339),
			TranscriptPath: transcriptPath,
		})
	}

	elapsed := time.Since(blockStart)
	if elapsed < 0 {
		elapsed = 0
	}

	switch item.BarMode {
	case "bar":
		return renderBar(elapsed, 5*time.Hour, 10)
	case "bar-short":
		return renderBar(elapsed, 5*time.Hour, 5)
	default:
		return util.FormatDuration(elapsed.Milliseconds())
	}
}

func renderBar(elapsed, total time.Duration, width int) string {
	filled := int(float64(elapsed) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
