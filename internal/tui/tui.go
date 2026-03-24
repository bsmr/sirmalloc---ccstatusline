package tui

import (
	"context"
	"fmt"
	"io"
)

func Run(_ context.Context, stdout, _ io.Writer) error {
	fmt.Fprintln(stdout, "ccstatusline-setup: TUI not yet implemented")
	return nil
}
