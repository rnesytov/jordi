package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	ResponseView struct {
		keyMap   ResponseKeyMap
		commands *Commands
		view     viewport.Model
		title    TitleView
		help     HelpView
	}
	ResponseKeyMap struct {
		resend key.Binding
	}
)

func DefaultResponseKeyMap() ResponseKeyMap {
	resend := key.NewBinding(key.WithKeys("ctrl+r"))
	resend.SetHelp(`ctrl+r`, "resend")

	return ResponseKeyMap{
		resend: resend,
	}
}

func (r ResponseKeyMap) Bindings() []key.Binding {
	return []key.Binding{r.resend}
}

func NewResponseView(commands *Commands) *ResponseView {
	view := viewport.New(0, 0)

	keyMap := DefaultResponseKeyMap()
	return &ResponseView{
		keyMap:   keyMap,
		commands: commands,
		view:     view,
		title:    NewTitleView("Response"),
		help:     NewHelpView(keyMap),
	}
}

func (r *ResponseView) Init() tea.Cmd {
	return nil
}

func (r *ResponseView) waitForMsg(sub <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (r *ResponseView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, r.keyMap.resend) {
			cmds = append(cmds, r.commands.ResendRequest())
		}
	case ShowResponseView:
		cmds = append(cmds, r.waitForMsg(msg.ch))
		cmds = append(cmds, r.commands.SetStatusLoading())
	case ReceivedResponse:
		r.view.SetContent(msg.Response)
		cmds = append(cmds, r.waitForMsg(msg.ch))
	case ReceivedStatus:
		statusMsgType := StatusMsgError
		if msg.Status == "OK" {
			statusMsgType = StatusMsgSuccess
		}
		cmds = append(cmds, func() tea.Msg {
			return NewStatusMessage{Msg: msg.Status, Type: statusMsgType}
		})
		cmds = append(cmds, r.commands.SetStatusOK())
	case Back:
		r.view.SetContent("")
		cmds = append(cmds, r.commands.ClearStatusMsg())
		cmds = append(cmds, r.commands.SetStatusOK())
	}
	var cmd tea.Cmd
	r.view, cmd = r.view.Update(msg)
	cmds = append(cmds, cmd)
	return r, tea.Batch(cmds...)
}

func (r *ResponseView) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, r.title.View(), r.view.View(), r.help.View())
}

func (r *ResponseView) HandleWindowSize(msg tea.WindowSizeMsg) {
	r.view.Width = msg.Width
	r.view.Height = msg.Height - helpHeight - titleHeight
	r.help.SetWidth(msg.Width)
}
