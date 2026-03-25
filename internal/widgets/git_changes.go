package widgets

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("git-changes", gitChangesWidget{})
}

type gitChangesWidget struct{}

func (gitChangesWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.CWD == "" {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "diff", "--numstat", "HEAD")
	cmd.Dir = si.CWD // Security R-7
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return ""
	}
	count := bytes.Count(out, []byte("\n"))
	// bytes.Count counts newlines; if output doesn't end with newline, add 1.
	if len(out) > 0 && out[len(out)-1] != '\n' {
		count++
	}
	if count == 1 {
		return fmt.Sprintf("%d file", count)
	}
	return fmt.Sprintf("%d files", count)
}
