package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"go.a8l.eu/ccstatusline/internal/pipe"
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

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return pipe.Run(ctx, stdout, stderr)
	}

	fmt.Fprintln(stderr, "use ccstatusline-setup for configuration")
	os.Exit(1)
	return nil
}
