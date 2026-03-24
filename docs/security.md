# Security Considerations — ccstatusline Go Port

This document captures security findings from a review of the TypeScript source and the Go architecture spec.
All findings must be addressed during Go implementation.

## Threat Model

ccstatusline runs as a local CLI tool invoked by Claude Code on every status update. The primary threat is not a remote attacker but rather:
- A compromised settings file (written by Claude's own tools or a malicious project)
- A manipulated stdin JSON payload (e.g., from a prompt-injected Claude session)
- Path traversal via user-controlled fields in the JSON input

## Implementation Rules (Go)

These rules are **mandatory** for the Go port:

### R-1: CustomCommand — no shell invocation

Use `exec.Command(path, args...)`, never `exec.Command("sh", "-c", command)`.
Shell features (pipes, redirects, globbing) are intentionally unsupported.

```go
// correct
cmd := exec.CommandContext(ctx, cfg.Command)

// forbidden
cmd := exec.CommandContext(ctx, "sh", "-c", cfg.Command)
```

Rationale: `cfg.Command` comes from `settings.json`, which could be written by Claude's tools. Any value in `command` would be executed as a shell script.

### R-2: Path validation for all external paths

All paths originating from stdin JSON or settings files must be validated before use:

```go
func validatePath(p, allowedBase string) (string, error) {
    clean := filepath.Clean(p)
    if !filepath.IsAbs(clean) {
        return "", fmt.Errorf("path must be absolute")
    }
    if !strings.HasPrefix(clean, allowedBase+string(filepath.Separator)) {
        return "", fmt.Errorf("path %q escapes allowed directory %q", clean, allowedBase)
    }
    return clean, nil
}
```

Apply to:
- `transcript_path` — allowed base: user's home directory
- `cwd` (used as `Cmd.Dir`) — allowed base: user's home directory
- Block timer cache path — allowed base: `~/.cache/ccstatusline/`

### R-3: session_id sanitization before use as filename

```go
var safeID = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func sanitizeSessionID(id string) (string, error) {
    if !safeID.MatchString(id) {
        return "", fmt.Errorf("session_id contains invalid characters")
    }
    return id, nil
}
```

### R-4: Atomic settings writes

Never use `os.WriteFile` for settings files. Always use temp-file + rename:

```go
func writeSettingsAtomic(path string, data []byte) error {
    dir := filepath.Dir(path)
    tmp, err := os.CreateTemp(dir, ".settings-*.json.tmp")
    if err != nil {
        return err
    }
    tmpName := tmp.Name()
    defer os.Remove(tmpName) // clean up on failure
    if _, err := tmp.Write(data); err != nil {
        tmp.Close()
        return err
    }
    if err := tmp.Close(); err != nil {
        return err
    }
    return os.Rename(tmpName, path)
}
```

Applies to: `~/.config/ccstatusline/settings.json` and `~/.claude/settings.json`.

### R-5: stdin size limit

```go
const maxStdinBytes = 1 << 20 // 1 MB
r := io.LimitReader(os.Stdin, maxStdinBytes)
if err := json.NewDecoder(r).Decode(&input); err != nil {
    return fmt.Errorf("decoding stdin: %w", err)
}
```

### R-6: CLAUDE_CONFIG_DIR validation

```go
func resolveClaudeConfigDir() (string, error) {
    dir := os.Getenv("CLAUDE_CONFIG_DIR")
    if dir == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        return filepath.Join(home, ".claude"), nil
    }
    clean := filepath.Clean(dir)
    if !filepath.IsAbs(clean) {
        return "", fmt.Errorf("CLAUDE_CONFIG_DIR must be an absolute path")
    }
    return clean, nil
}
```

### R-7: git Cmd.Dir must be set and validated

```go
func runGit(ctx context.Context, cwd string, args ...string) (string, error) {
    // cwd must be validated via R-2 before calling this function
    cmd := exec.CommandContext(ctx, "git", args...)
    cmd.Dir = cwd
    out, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}
```

## Findings Reference

| ID | Severity | Source | Description |
|----|----------|--------|-------------|
| GC-1 | Critical | Go spec | CustomCommand with `sh -c` → shell injection → addressed by R-1 |
| GC-2 | Critical | Go spec | `transcript_path` without validation → path traversal → addressed by R-2 |
| H-3/GH-2 | High | TS + Go | `session_id` as filename without sanitization → addressed by R-3 |
| H-1 | High | TS | `transcript_path` read without validation → addressed by R-2 |
| H-2 | High | TS | Subagent glob based on manipulable `transcript_path` → addressed by R-2 |
| GH-1 | High | Go spec | `gitBranch` snippet missing `Cmd.Dir` → addressed by R-7 |
| GM-2 | Medium | Go spec | Non-atomic settings write → addressed by R-4 |
| GM-3/M-1 | Medium | TS + Go | `CLAUDE_CONFIG_DIR` without restriction → addressed by R-6 |
| GM-4 | Low | Go spec | No stdin size limit → addressed by R-5 |
| M-2 | Medium | TS | `claude-settings.ts` parses JSON without schema validation (TypeScript only) |
