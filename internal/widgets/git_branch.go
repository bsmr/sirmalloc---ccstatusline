package widgets

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("git-branch", gitBranchWidget{})
}

type gitBranchWidget struct{}

func (gitBranchWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil || si.CWD == "" {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")
	cmd.Dir = si.CWD // Security R-7: set working directory from input
	out, err := cmd.Output()
	if err != nil {
		// not a git repo or git not found — degrade gracefully
		return ""
	}
	branch := strings.TrimSpace(string(out))
	if item.HideNoGit && branch == "" {
		return ""
	}
	return branch
}
