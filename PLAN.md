# Development Plan — ccstatusline Go Port

## Phase 1: Minimal Piped Mode (renderer) ✅ DONE

1. `cmd/ccstatusline/main.go` — isatty check, `main() → run()` pattern, stdin read with size limit
2. `internal/input/parse.go` — JSON decode → `StatusInput`
3. `internal/render/ansi.go` — basic color codes
4. First widgets: `model`, `context-percentage`, `session-clock`, `session-cost`, `separator`, `flex-separator`
5. `internal/render/renderer.go` — assemble widgets, write to stdout
6. `internal/config/settings.go` — `Settings`/`WidgetItem` structs, `Load()`/`Save()` (atomic)
7. Smoke test: `cat testdata/example_input.json | go run ./cmd/ccstatusline`

## Phase 1.5: TUI Skeleton — view & edit existing widgets ✅ DONE

Goal: after this phase `ccstatusline-setup` shows the current config and allows
editing all Phase-1 widgets. No new renderer features.

8. Add Bubble Tea + Lip Gloss to `go.mod` (`cmd/ccstatusline-setup` only) ✅
9. Main screen: list of lines ([][]WidgetItem), navigate with j/k ✅
10. Line screen: list of widgets in the selected line, navigate with j/k ✅
11. Widget screen: edit color/fg/bg/bold for the selected widget ✅
12. Widget-type-specific options for Phase-1 widgets ✅
13. Add widget (type selection list), remove widget (d) ✅
14. Save with Ctrl+S (atomic write via `config.Save()`) ✅
15. Quit with q/Esc from top level ✅

## Phase 2: All Widgets (renderer + TUI editor page per widget) ✅ DONE

All widgets implemented. TUI supports all widget types via `fieldsFor()`.

16. `git-branch` — `os/exec`, 2s timeout, `Cmd.Dir` validation ✅
17. `git-worktree` — same exec pattern ✅
18. `cwd` — path display, fish-style, segments ✅
19. `version` — static from `StatusInput` ✅
20. `output-style` — pass-through ✅
21. `tokens-in`, `tokens-out`, `tokens-total` — from `StatusInput` ✅
22. `context-length` — from `StatusInput.Model` ✅
23. `context-pct-usable` — derived from tokens + model limit ✅
24. `rate-limit-five-hour`, `rate-limit-seven-day` — from `StatusInput.RateLimits` ✅
25. `vim-mode` — from `StatusInput.Vim` ✅
26. `terminal-width` — from termWidth parameter ✅
27. `custom-text` — static, `item.Text` ✅
28. `custom-command` — `exec.Command(path, args...)`, no `sh -c` (R-1), timeout ✅
29. `block-timer` — JSONL parsing, cache in `~/.cache/ccstatusline/` ✅

## Phase 3: Settings & Powerline (renderer + TUI screens) ✅ DONE

30. Powerline rendering, separators, caps (renderer) ✅
31. Flex separator, ANSI/OSC-aware truncation (renderer) ✅
32. Terminal width modes (`full`, `full-minus-40`, `full-until-compact`) ✅
33. TUI: Powerline setup screen (enabled, separator, caps, autoAlign, theme) ✅
34. TUI: Global options screen (flexMode, colorLevel, compactThreshold, globalBold) ✅
35. Powerline theme system: 11 built-in themes × 3 color levels ✅

## Phase 4: TUI Polish & Integration ✅ DONE

36. Live preview (render current config with synthetic input) ✅
37. Install/Uninstall in Claude Code settings (`~/.claude/settings.json`) ✅
38. Color picker: Basic color selection (8 named colors) ✅

## Phase 5: Polish

38. Windows compatibility (paths, code page)
39. Full test coverage
40. README
41. Goreleaser for automated releases

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
