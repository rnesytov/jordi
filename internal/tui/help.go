package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

const (
	helpHeight = 1
)

var (
	helpStyle = lipgloss.NewStyle().PaddingLeft(2)
)

type (
	BindingsGetter interface {
		Bindings() []key.Binding
	}
	HelpView struct {
		bindings BindingsGetter
		view     help.Model
	}
)

func NewHelpView(bindings BindingsGetter) HelpView {
	return HelpView{bindings: bindings, view: help.NewModel()}
}

func (h *HelpView) View() string {
	return helpStyle.Render(h.view.ShortHelpView(h.bindings.Bindings()))
}

func (h *HelpView) SetWidth(width int) {
	h.view.Width = width
}
