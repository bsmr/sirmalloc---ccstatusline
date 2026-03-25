# ccstatusline

A customizable status line renderer for the Claude Code CLI.

## Overview

ccstatusline replaces the built-in Claude Code status line with a configurable bar
that displays model information, context usage, session cost, elapsed time, git branch,
and more. It is designed to start fast (< 10 ms) and have minimal dependencies in the
renderer path.

The project consists of two binaries:

- `ccstatusline` — lean renderer; reads JSON from stdin, writes an ANSI status line to
  stdout. Only dependency: `golang.org/x/term`.
- `ccstatusline-setup` — interactive TUI editor for configuration (Bubble Tea + Lip Gloss).

## Installation

### Build from source

```bash
go build -o bin/ccstatusline       ./cmd/ccstatusline
go build -o bin/ccstatusline-setup ./cmd/ccstatusline-setup
```

Cross-compilation examples:

```bash
GOOS=linux  GOARCH=amd64 go build -o bin/dist/ccstatusline-linux-amd64  ./cmd/ccstatusline
GOOS=darwin GOARCH=arm64 go build -o bin/dist/ccstatusline-darwin-arm64 ./cmd/ccstatusline
```

### Download a release binary

Pre-built binaries for Linux, macOS, and Windows are available on the
[Releases](../../releases) page once a tagged release has been published.

### Integrate with Claude Code

**Option A — via setup tool (recommended):**

```bash
bin/ccstatusline-setup
# Press `i` to install into ~/.claude/settings.json automatically.
```

**Option B — manual:**

Add the following to `~/.claude/settings.json` (or `$CLAUDE_CONFIG_DIR/settings.json`):

```json
{
  "statusLine": {
    "type": "command",
    "command": "/absolute/path/to/ccstatusline",
    "padding": 0
  }
}
```

## Configuration

Run `ccstatusline-setup` to open the interactive TUI editor. Settings are stored in
`~/.config/ccstatusline/settings.json` and are JSON-compatible with the TypeScript
original (settings version 3).

### TUI keyboard shortcuts

| Key | Action |
|-----|--------|
| `j` / `k` / arrow keys | Navigate |
| `Enter` | Select / confirm |
| `Esc` / `q` | Back / cancel |
| `a` | Add widget |
| `d` | Remove widget |
| `Ctrl+S` | Save |
| `p` | Toggle bar mode (Block Timer) |
| `h` | Toggle hide-no-git |
| `l` / `u` | Toggle remaining mode (Context %) |
| `i` | Install in Claude Code settings |
| `u` | Uninstall from Claude Code settings |

## Widget Reference

| Widget type | Description |
|-------------|-------------|
| `model` | Model display name from Claude Code |
| `git-branch` | Current git branch (`git branch --show-current`) |
| `git-changes` | Line insertions + deletions from `git diff --numstat` |
| `git-worktree` | Active git worktree name |
| `session-clock` | Total session duration, e.g. `2hr 15m` |
| `session-cost` | Total session cost, e.g. `$1.23` |
| `block-timer` | Elapsed time within the current 5-hour usage block |
| `cwd` | Current working directory (configurable segments, fish-style abbreviation) |
| `version` | Claude Code version string |
| `output-style` | Output style name |
| `tokens-input` | Total input tokens, e.g. `15.2k` |
| `tokens-output` | Total output tokens |
| `tokens-cached` | Cache-read input tokens |
| `tokens-total` | Combined input + output token count |
| `context-length` | Context window size |
| `context-percentage` | Context window usage percentage (toggle: used / remaining) |
| `context-pct-usable` | Percentage relative to auto-compact boundary (80% of max) |
| `rate-limit-five-hour` | 5-hour rate limit used percentage |
| `rate-limit-seven-day` | 7-day rate limit used percentage |
| `vim-mode` | Vim mode indicator (only shown when active) |
| `terminal-width` | Current terminal width (debugging) |
| `custom-text` | Static text from config (emoji supported) |
| `custom-command` | Output of an external command (no shell; direct exec only) |
| `separator` | Configurable separator character (`|`, `-`, `,`, space) |
| `flex-separator` | Fills available space; used for right-alignment |

## Architecture

```
cmd/ccstatusline/       renderer — isatty check, stdin JSON decode, ANSI output
cmd/ccstatusline-setup/ TUI editor — Bubble Tea app, reads/writes settings.json
internal/config/        settings load/save, Claude Code integration, XDG paths
internal/input/         JSON stdin -> StatusInput struct
internal/render/        renderer, ANSI helpers, powerline, truncation
internal/widgets/       one file per widget type
internal/tui/           TUI screens (setup binary only)
internal/util/          formatting helpers, git utilities, terminal detection
```

`cmd/ccstatusline` has no TUI dependencies. Mode is determined by `term.IsTerminal`:
if stdin is not a terminal (piped), the renderer runs; otherwise an error message
directs the user to `ccstatusline-setup`.

## Security

Custom commands are executed via `exec.Command(path, args...)` — never via `sh -c`.
Shell features (pipes, redirects, variable expansion) are intentionally unsupported.
All paths from external sources are validated against their allowed base directories.
Stdin is wrapped with a 1 MiB size limit before decoding.

See [docs/security.md](docs/security.md) for the full list of security rules and
their implementation.

## Development

```bash
# Run all tests
go test ./...

# Vet
go vet ./...

# Lint (requires golangci-lint)
golangci-lint run

# Smoke test — piped mode
cat testdata/example_input.json | go run ./cmd/ccstatusline
```

## License

MIT License — see LICENSE file for details.
