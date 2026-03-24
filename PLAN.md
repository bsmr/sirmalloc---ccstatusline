# Development Plan — ccstatusline Go Port

## Phase 1: Minimal Piped Mode (renderer, Priority: HIGH)

1. `cmd/ccstatusline/main.go` — isatty check, `main() → run()` pattern, stdin read with size limit
2. `internal/input/parse.go` — JSON decode → `StatusInput`
3. `internal/render/ansi.go` — basic color codes
4. First widgets: `model`, `contextPct`, `sessionClock`
5. `internal/render/renderer.go` — assemble widgets, write to stdout
6. Smoke test: `cat testdata/example_input.json | go run ./cmd/ccstatusline`

## Phase 2: All Widgets

7. All widgets from the table in `docs/spec.md`
8. Git widgets with `os/exec`, 2s timeout, `Cmd.Dir` validation (see `docs/security.md` R-7)
9. Block timer with JSONL parsing and cache
10. Custom command with stdin forwarding (no `sh -c` — see `docs/security.md` R-1)
11. Rate limits widgets (`rateLimitFiveHour`, `rateLimitSevenDay`)
12. Vim mode widget

## Phase 3: Settings & Powerline

13. Settings file read/write with atomic writes (see `docs/security.md` R-4)
14. Powerline rendering, separators, caps
15. Flex separator, auto-alignment
16. Terminal width modes
17. ANSI/OSC-aware truncation

## Phase 4: TUI (ccstatusline-setup binary, Priority: MEDIUM)

18. Separate `cmd/ccstatusline-setup` entry point
19. Bubble Tea skeleton
20. Main menu, Line selector
21. Items editor with widget selection
22. Widget editor with color/option configuration
23. Color menu (Basic / 256 / Truecolor)
24. Powerline setup screen
25. Global options screen
26. Live preview
27. Install/Uninstall in Claude Code settings

## Phase 5: Polish

28. Windows compatibility (paths, code page)
29. Full test coverage
30. README
31. Goreleaser for automated releases

## Quality Criteria

- [ ] `go vet ./...` clean
- [ ] `golangci-lint run` clean
- [ ] All widgets render correctly with `testdata/example_input.json`
- [ ] nil fields cause no panics
- [ ] Git widgets degrade gracefully outside a git repo
- [ ] Settings file is JSON-compatible with the TypeScript original
- [ ] Binary starts in < 10ms (vs ~200ms+ for npx/bunx)
- [ ] All security rules in `docs/security.md` implemented and tested
- [ ] Cross-compilation works for Linux, macOS, Windows
- [ ] `ccstatusline` binary has no TUI dependencies (verify with `go mod graph`)

## TUI Keyboard Shortcuts (from TypeScript original)

| Key | Action |
|-----|--------|
| `Enter` | Select / confirm |
| `Esc` / `q` | Back / cancel |
| `j/k` / `↑/↓` | Navigate |
| `a` | Add widget |
| `d` | Remove widget |
| `Ctrl+S` | Save |
| `p` | Toggle bar mode (Block Timer) |
| `h` | Toggle hide-no-git |
| `l/u` | Toggle remaining mode (Context %) |
| `i` | Install in Claude Code settings |
| `u` | Uninstall from Claude Code settings |
