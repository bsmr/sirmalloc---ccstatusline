package pipe

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.a8l.eu/ccstatusline/internal/input"
)

func Run(_ context.Context, stdout, _ io.Writer) error {
	si, err := input.Decode(io.LimitReader(os.Stdin, 1<<20))
	if err != nil {
		return err
	}
	_ = si
	fmt.Fprintln(stdout, "ccstatusline: not yet implemented")
	return nil
}
