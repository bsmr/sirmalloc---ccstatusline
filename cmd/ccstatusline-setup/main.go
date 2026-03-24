package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"go.a8l.eu/ccstatusline/internal/tui"
)

func main() {
	if err := run(context.Background(), os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	return tui.Run(ctx, stdout, stderr)
}
