# Technical Specification — ccstatusline Go Port

## JSON Input Format (stdin from Claude Code)

Claude Code sends the following JSON on stdin with every status update.
All fields are optional/nullable — especially before the first API call.

```json
{
  "hook_event_name": "Status",
  "session_id": "abc123...",
  "transcript_path": "/path/to/transcript.jsonl",
  "cwd": "/current/working/directory",
  "model": {
    "id": "claude-opus-4-6[1m]",
    "display_name": "Opus 4.6 (1M context)"
  },
  "workspace": {
    "current_dir": "/current/working/directory",
    "project_dir": "/original/project/directory",
    "added_dirs": []
  },
  "version": "2.1.80",
  "output_style": { "name": "default" },
  "cost": {
    "total_cost_usd": 0.01234,
    "total_duration_ms": 45000,
    "total_api_duration_ms": 2300,
    "total_lines_added": 156,
    "total_lines_removed": 23
  },
  "context_window": {
    "total_input_tokens": 50113,
    "total_output_tokens": 10462,
    "context_window_size": 1000000,
    "used_percentage": 8,
    "remaining_percentage": 92,
    "current_usage": {
      "input_tokens": 8500,
      "output_tokens": 1200,
      "cache_creation_input_tokens": 5000,
      "cache_read_input_tokens": 2000
    }
  },
  "exceeds_200k_tokens": false,
  "rate_limits": {
    "five_hour": { "used_percentage": 42, "resets_at": 1774020000 },
    "seven_day":  { "used_percentage": 15, "resets_at": 1774540000 }
  },
  "vim": { "mode": "NORMAL" }
}
```

## Go Types

```go
type StatusInput struct {
    HookEventName  string          `json:"hook_event_name"`
    SessionID      string          `json:"session_id"`
    TranscriptPath string          `json:"transcript_path"`
    CWD            string          `json:"cwd"`
    Model          *ModelInfo      `json:"model"`
    Workspace      *WorkspaceInfo  `json:"workspace"`
    Version        string          `json:"version"`
    OutputStyle    *OutputStyle    `json:"output_style"`
    Cost           *CostInfo       `json:"cost"`
    ContextWindow  *ContextWindow  `json:"context_window"`
    Exceeds200k    bool            `json:"exceeds_200k_tokens"`
    RateLimits     *RateLimits     `json:"rate_limits"`
    Vim            *VimInfo        `json:"vim"`
}

type ModelInfo struct {
    ID          string `json:"id"`
    DisplayName string `json:"display_name"`
}

type WorkspaceInfo struct {
    CurrentDir string   `json:"current_dir"`
    ProjectDir string   `json:"project_dir"`
    AddedDirs  []string `json:"added_dirs"`
}

type OutputStyle struct {
    Name string `json:"name"`
}

type CostInfo struct {
    TotalCostUSD        float64 `json:"total_cost_usd"`
    TotalDurationMs     int64   `json:"total_duration_ms"`
    TotalAPIDurationMs  int64   `json:"total_api_duration_ms"`
    TotalLinesAdded     int     `json:"total_lines_added"`
    TotalLinesRemoved   int     `json:"total_lines_removed"`
}

type ContextWindow struct {
    TotalInputTokens  int           `json:"total_input_tokens"`
    TotalOutputTokens int           `json:"total_output_tokens"`
    ContextWindowSize int           `json:"context_window_size"`
    UsedPercentage    float64       `json:"used_percentage"`
    RemainingPct      float64       `json:"remaining_percentage"`
    CurrentUsage      *CurrentUsage `json:"current_usage"`
}

type CurrentUsage struct {
    InputTokens              int `json:"input_tokens"`
    OutputTokens             int `json:"output_tokens"`
    CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
    CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

type RateLimits struct {
    FiveHour  *RateLimitPeriod `json:"five_hour"`
    SevenDay  *RateLimitPeriod `json:"seven_day"`
}

type RateLimitPeriod struct {
    UsedPercentage *float64 `json:"used_percentage"` // nullable
    ResetsAt       *int64   `json:"resets_at"`        // Unix epoch seconds, nullable
}

type VimInfo struct {
    Mode string `json:"vim_mode"`
}
```

## Widget Interface

```go
type Widget interface {
    Render(input *StatusInput, cfg WidgetConfig) string
}
```

## Widget List

| Widget type | Source | Notes |
|-------------|--------|-------|
| `model` | `input.Model.DisplayName` | RawValue: value only, no label |
| `gitBranch` | `git branch --show-current` | via `os/exec`, 2s timeout, `Cmd.Dir = cwd` |
| `gitChanges` | `git diff --numstat` | counts insertions + deletions |
| `gitWorktree` | `git worktree list` | active worktree name |
| `sessionClock` | `input.Cost.TotalDurationMs` | format: `2hr 15m` |
| `sessionCost` | `input.Cost.TotalCostUSD` | format: `$1.23` |
| `blockTimer` | transcript JSONL parse | elapsed time in 5h block |
| `cwd` | `input.Workspace.CurrentDir` | configurable segment count, fish-style abbreviation |
| `version` | `input.Version` | |
| `outputStyle` | `input.OutputStyle.Name` | |
| `tokensInput` | `input.ContextWindow.TotalInputTokens` | formatted: `15.2k` |
| `tokensOutput` | `input.ContextWindow.TotalOutputTokens` | |
| `tokensCached` | `CurrentUsage.CacheReadInputTokens` | |
| `tokensTotal` | input + output | |
| `contextLength` | `input.ContextWindow.ContextWindowSize` | |
| `contextPct` | `input.ContextWindow.UsedPercentage` | toggle: used/remaining |
| `contextPctUsable` | 80% of max (auto-compact boundary) | 1M models → 800k, else 160k |
| `rateLimitFiveHour` | `input.RateLimits.FiveHour.UsedPercentage` | shows resets_at on hover if supported |
| `rateLimitSevenDay` | `input.RateLimits.SevenDay.UsedPercentage` | |
| `vimMode` | `input.Vim.Mode` | only rendered when non-empty |
| `terminalWidth` | `golang.org/x/term` | debugging widget |
| `customText` | static text from config | emoji supported |
| `customCommand` | shell command (no `sh -c`) | JSON forwarded to stdin |
| `separator` | configurable character | `\|`, `-`, `,`, space |
| `flexSeparator` | fills available space | used for right-alignment |

## Model Context Detection

```go
func getModelContext(modelID string) (maxTokens int, usableTokens int) {
    if strings.Contains(strings.ToLower(modelID), "[1m]") {
        return 1_000_000, 800_000
    }
    return 200_000, 160_000
}
```

## Formatting Helpers

```go
func formatTokens(n int) string {
    switch {
    case n >= 1_000_000:
        return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
    case n >= 1_000:
        return fmt.Sprintf("%.1fk", float64(n)/1_000)
    default:
        return strconv.Itoa(n)
    }
}

func formatDuration(ms int64) string {
    d := time.Duration(ms) * time.Millisecond
    h := int(d.Hours())
    m := int(d.Minutes()) % 60
    if h > 0 {
        return fmt.Sprintf("%dhr %dm", h, m)
    }
    return fmt.Sprintf("%dm", m)
}
```

## Settings Schema

Path: `~/.config/ccstatusline/settings.json`

Must remain JSON-compatible with the TypeScript original for seamless migration.

```go
type Settings struct {
    Lines         []LineConfig     `json:"lines"`
    GlobalOptions GlobalOptions    `json:"globalOptions"`
    Powerline     *PowerlineConfig `json:"powerline,omitempty"`
    TerminalWidth string           `json:"terminalWidth"` // "full", "full-minus-40", "full-until-compact"
    ColorMode     string           `json:"colorMode"`     // "basic", "256", "truecolor"
}

type LineConfig struct {
    Widgets []WidgetConfig `json:"widgets"`
}

type WidgetConfig struct {
    Type      string `json:"type"`
    FG        string `json:"fg"`
    BG        string `json:"bg"`
    Bold      bool   `json:"bold"`
    RawValue  bool   `json:"rawValue"`
    Merge     bool   `json:"merge"`
    // widget-specific:
    Text      string `json:"text,omitempty"`
    Command   string `json:"command,omitempty"`
    Timeout   int    `json:"timeout,omitempty"`   // ms
    Segments  int    `json:"segments,omitempty"`
    FishStyle bool   `json:"fishStyle,omitempty"`
    HideNoGit bool   `json:"hideNoGit,omitempty"`
    Remaining bool   `json:"remaining,omitempty"`
    BarMode   string `json:"barMode,omitempty"`   // "time", "bar", "bar-short"
    SepChar   string `json:"sepChar,omitempty"`
}

type GlobalOptions struct {
    DefaultPadding   int    `json:"defaultPadding"`
    DefaultSeparator string `json:"defaultSeparator"`
    InheritColors    bool   `json:"inheritColors"`
    GlobalBold       bool   `json:"globalBold"`
    OverrideFG       string `json:"overrideFG,omitempty"`
    OverrideBG       string `json:"overrideBG,omitempty"`
}

type PowerlineConfig struct {
    Enabled   bool   `json:"enabled"`
    SepChar   string `json:"sepChar"`   // default: U+E0B0
    CapLeft   string `json:"capLeft"`
    CapRight  string `json:"capRight"`
    AutoAlign bool   `json:"autoAlign"`
}
```

## Color System

Three modes detected from `$COLORTERM` / `$TERM`:

| Mode | Escape sequence | Detection |
|------|----------------|-----------|
| Basic (16) | `\033[3Xm` / `\033[4Xm` | fallback |
| 256 | `\033[38;5;{n}m` / `\033[48;5;{n}m` | `$TERM=xterm-256color` |
| Truecolor | `\033[38;2;{r};{g};{b}m` | `$COLORTERM=truecolor` |

Recommended library: `github.com/muesli/termenv` (auto-detects capabilities).

## Powerline Mode

- Arrow separators: `` (U+E0B0), `` (U+E0B2)
- Each widget has its own FG/BG color
- Transitions: previous widget's BG becomes the separator's FG
- Auto-alignment across lines (equal column widths)
- Widget merging: shared background for adjacent widgets

## Output Format

- One or more lines on stdout, separated by `\n`
- ANSI/OSC-aware truncation when line exceeds terminal width
- Ellipsis (`...`) appended when truncated
- OSC 8 hyperlink support for link-capable widgets

## Terminal Width Modes

| Mode | Behavior |
|------|---------|
| `full` | Always use full terminal width |
| `full-minus-40` | Reserve 40 chars for auto-compact message (default) |
| `full-until-compact` | Dynamic, based on context-% threshold |

## Block Timer

Reads `transcript_path` from the JSON, parses the JSONL file, finds the timestamp of the current 5-hour block (rounded down to the hour), computes elapsed time.
Cache: `~/.cache/ccstatusline/block-cache-<sha256(configDir)>.json`

## Project Structure

```
ccstatusline/
├── CLAUDE.md
├── PLAN.md
├── go.mod
├── go.sum
├── testdata/
│   └── example_input.json
├── docs/
│   ├── spec.md               (this file)
│   ├── security.md
│   └── typescript-reference.md
├── cmd/
│   ├── ccstatusline/
│   │   └── main.go           entry point, mode detection, calls internal/pipe
│   └── ccstatusline-setup/
│       └── main.go           entry point, calls internal/tui
├── internal/
│   ├── config/
│   │   ├── settings.go       read/write settings.json
│   │   ├── claude.go         Claude Code settings.json integration
│   │   └── paths.go          config paths, XDG, CLAUDE_CONFIG_DIR
│   ├── input/
│   │   └── parse.go          JSON stdin → StatusInput
│   ├── render/
│   │   ├── renderer.go       main renderer: settings + input → lines
│   │   ├── ansi.go           ANSI escape codes, colors
│   │   ├── powerline.go      powerline rendering, separators, caps
│   │   └── truncate.go       ANSI/OSC-aware truncation
│   ├── widgets/
│   │   ├── widget.go         Widget interface + registry
│   │   ├── model.go
│   │   ├── git_branch.go
│   │   ├── git_changes.go
│   │   ├── git_worktree.go
│   │   ├── session_clock.go
│   │   ├── session_cost.go
│   │   ├── block_timer.go
│   │   ├── cwd.go
│   │   ├── version.go
│   │   ├── output_style.go
│   │   ├── tokens.go
│   │   ├── context.go
│   │   ├── rate_limits.go
│   │   ├── vim_mode.go
│   │   ├── terminal_width.go
│   │   ├── custom_text.go
│   │   ├── custom_command.go
│   │   ├── separator.go
│   │   └── flex_separator.go
│   ├── tui/                  (ccstatusline-setup only)
│   │   ├── app.go
│   │   ├── mainmenu.go
│   │   ├── line_selector.go
│   │   ├── items_editor.go
│   │   ├── widget_editor.go
│   │   ├── color_menu.go
│   │   ├── powerline_setup.go
│   │   ├── global_options.go
│   │   └── preview.go
│   └── util/
│       ├── format.go         token/duration formatting
│       ├── git.go            git helpers
│       └── term.go           terminal width, color mode detection
└── README.md
```

## Dependencies

| Module | Used by | Purpose |
|--------|---------|---------|
| `golang.org/x/term` | ccstatusline | isatty, terminal width |
| `github.com/charmbracelet/bubbletea` | ccstatusline-setup | TUI framework |
| `github.com/charmbracelet/lipgloss` | ccstatusline-setup | TUI styling |
| `github.com/charmbracelet/bubbles` | ccstatusline-setup | TUI components |
| `github.com/muesli/termenv` | ccstatusline | color profile detection |

## Claude Code Integration

```json
{
  "statusLine": {
    "type": "command",
    "command": "/path/to/bin/ccstatusline",
    "padding": 0
  }
}
```
