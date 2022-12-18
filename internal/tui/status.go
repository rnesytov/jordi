package tui

import tea "github.com/charmbracelet/bubbletea"

type (
	StatusBar struct {
		currentStatus string
		Error         string
	}
)

func NewStatusBar() *StatusBar {
	return &StatusBar{}
}

func (s *StatusBar) Init() tea.Cmd {
	return nil
}

func (s *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ShowStatus:
		s.currentStatus = msg.Status
	}
	return s, nil
}

func (s *StatusBar) View() string {
	return s.currentStatus
}
