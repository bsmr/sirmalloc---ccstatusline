package pipe

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/render"
)

func Run(_ context.Context, stdout, _ io.Writer) error {
	si, err := input.Decode(io.LimitReader(os.Stdin, 1<<20))
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if termWidth <= 0 {
		termWidth = 80
	}

	lines := render.Render(si, cfg, termWidth)
	fmt.Fprintln(stdout, strings.Join(lines, "\n"))
	return nil
}
