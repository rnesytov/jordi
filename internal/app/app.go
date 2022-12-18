package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/profx5/jordi/internal/config"
	"github.com/profx5/jordi/internal/grpc"
	"github.com/profx5/jordi/internal/tui"
)

type App struct {
	config config.Config
}

func New(config config.Config) *App {
	return &App{config: config}
}

func (a *App) Run(ctx context.Context) error {
	gw, err := grpc.New(ctx, a.config.Target, grpc.DefaultOpts())
	if err != nil {
		return err
	}
	defer gw.Close()

	root := tui.NewRoot(gw)

	p := tea.NewProgram(root, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}