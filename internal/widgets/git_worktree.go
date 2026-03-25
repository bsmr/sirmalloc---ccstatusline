package widgets

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("git-worktree", gitWorktreeWidget{})
}

type gitWorktreeWidget struct{}

func (gitWorktreeWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil || si.CWD == "" {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Use git rev-parse --show-toplevel to get the worktree root path,
	// then return its base name as the worktree name.
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	cmd.Dir = si.CWD // Security R-7
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	toplevel := strings.TrimSpace(string(out))
	if toplevel == "" {
		return ""
	}
	name := filepath.Base(toplevel)
	if item.HideNoGit && name == "" {
		return ""
	}
	return name
}
