package tui

import (
	"context"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"go.a8l.eu/ccstatusline/internal/config"
)

func Run(_ context.Context, _, _ io.Writer) error {
	settings, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	path, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("resolving config path: %w", err)
	}

	m := newModel(settings, path)
	p := tea.NewProgram(m, tea.WithOutput(os.Stdout))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tui: %w", err)
	}
	return nil
}
