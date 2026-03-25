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
	Register("custom-command", customCommandWidget{})
}

type customCommandWidget struct{}

func (customCommandWidget) Render(_ *input.StatusInput, item config.WidgetItem) string {
	if item.Command == "" {
		return ""
	}
	timeout := time.Duration(item.Timeout) * time.Millisecond
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Security R-1: use exec.Command(path, args...), never sh -c.
	cmd := exec.CommandContext(ctx, item.Command)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
