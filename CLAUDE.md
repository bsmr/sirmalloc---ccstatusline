# CLAUDE.md — ccstatusline Go Port

## Quick Start

```bash
# Build (binaries always go to bin/)
go build -o bin/ccstatusline ./cmd/ccstatusline
go build -o bin/ccstatusline-setup ./cmd/ccstatusline-setup

# Cross-compilation
GOOS=linux   GOARCH=amd64 go build -o bin/dist/ccstatusline-linux-amd64 ./cmd/ccstatusline
GOOS=darwin  GOARCH=arm64 go build -o bin/dist/ccstatusline-darwin-arm64 ./cmd/ccstatusline

# Test
go test ./...
go vet ./...
golangci-lint run

# Test piped mode
cat testdata/example_input.json | go run ./cmd/ccstatusline
```

## Architecture

Two binaries with a strict separation of concerns:

| Binary | Deps | Purpose |
|--------|------|---------|
| `cmd/ccstatusline` | `golang.org/x/term` only | Renderer: reads JSON from stdin, writes status line to stdout |
| `cmd/ccstatusline-setup` | Bubble Tea + Lip Gloss | Interactive TUI: reads/writes `~/.config/ccstatusline/settings.json` |

Mode detection in `cmd/ccstatusline`:
```go
if !term.IsTerminal(int(os.Stdin.Fd())) {
    // piped mode — render status line
} else {
    fmt.Fprintln(os.Stderr, "use ccstatusline-setup for configuration")
    os.Exit(1)
}
```

All logic lives in `run()`, never in `main()`. The `main() → run()` pattern is mandatory — see global CLAUDE.md.

Settings path: `~/.config/ccstatusline/settings.json` (shared format with TypeScript original).
Claude Code integration: `~/.claude/settings.json` (or `$CLAUDE_CONFIG_DIR/settings.json`).

## Key Design Decisions

**CustomCommand widget:** Uses `exec.Command(path, args...)` — never `sh -c`. Shell features (pipes, redirects) are intentionally unsupported. This is a deliberate security trade-off.

**Path validation:** All paths from external sources (stdin JSON, settings) must be validated:
```go
clean := filepath.Clean(p)
if !strings.HasPrefix(clean, allowedBase+string(filepath.Separator)) {
    return "", fmt.Errorf("path escapes allowed directory")
}
```
Applies to: `transcript_path`, `cwd` (as `Cmd.Dir`), `session_id` (as filename component).

**`session_id` as filename:** Sanitize to `[a-zA-Z0-9_-]` only before using as a filename part.

**Atomic settings writes:** Use `os.CreateTemp` + `os.Rename` — never `os.WriteFile` directly.

**stdin size limit:** Wrap with `io.LimitReader(os.Stdin, 1<<20)` before decoding.

## Gotchas

- All JSON fields can be `nil` — pointer types everywhere, nil-checks before use
- ANSI-aware string length: strip `\033[...m` and OSC-8 sequences before measuring visible width
- `git` commands need `Cmd.Dir` set to `cwd` from the JSON input; use `exec.LookPath` to detect absence
- Block timer cache lives in `~/.cache/ccstatusline/block-cache-<hash>.json`
- `bin/` must be in `.gitignore`

## References

- [Full technical spec](docs/spec.md) — JSON input format, Go types, widget list, settings schema, color system
- [Security considerations](docs/security.md) — findings and implementation rules for Go
- [Development plan](PLAN.md) — phases, roadmap, quality criteria
- [TypeScript reference](docs/typescript-reference.md) — original implementation (read-only reference)
