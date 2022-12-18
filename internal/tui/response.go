package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	ResponseView struct {
		view viewport.Model
	}
)

func NewResponseView() *ResponseView {
	return &ResponseView{
		view: viewport.New(0, 0),
	}
}

func (r *ResponseView) Init() tea.Cmd {
	return nil
}

func (r *ResponseView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case ShowResponse:
		r.view.SetContent(msg.Response)
	}
	var cmd tea.Cmd
	r.view, cmd = r.view.Update(msg)
	cmds = append(cmds, cmd)
	return r, tea.Batch(cmds...)
}

func (r *ResponseView) View() string {
	return r.view.View()
}

func (m *ResponseView) HandleWindowSize(msg tea.WindowSizeMsg) {
	m.view.Width = msg.Width
	m.view.Height = msg.Height
}
