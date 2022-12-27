package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusMsgType int

const (
	StatusMsgError   StatusMsgType = iota
	StatusMsgSuccess StatusMsgType = iota
)

var (
	statusBackgorundColor = lipgloss.Color("#555753")
	statusBarStyle        = lipgloss.NewStyle().
				Background(statusBackgorundColor)
	statusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#4e9a06")).
			PaddingLeft(2).
			PaddingRight(2).
			Bold(true)
	msgStylesMap = map[StatusMsgType]lipgloss.Style{
		StatusMsgSuccess: lipgloss.NewStyle().
			Background(statusBackgorundColor).
			Foreground(lipgloss.Color("#16a402")).
			PaddingLeft(2).
			PaddingRight(2),
		StatusMsgError: lipgloss.NewStyle().
			Background(statusBackgorundColor).
			Foreground(lipgloss.Color("#cc0000")).
			PaddingLeft(2).
			PaddingRight(2),
	}
)

type (
	StatusView struct {
		currentStatus string
		msg           string
		statusMsgType StatusMsgType

		width int
	}
)

func NewStatusView() *StatusView {
	return &StatusView{}
}

func (s *StatusView) Init() tea.Cmd {
	return nil
}

func (s *StatusView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetStatus:
		s.currentStatus = msg.Status
	case SetStatusMsg:
		s.msg = msg.Msg
		s.statusMsgType = msg.Type
	case ClearStatusMsg:
		s.msg = ""
	}
	return s, nil
}

func (s *StatusView) View() string {
	views := []string{statusStyle.Render(s.currentStatus)}
	if s.msg != "" {
		views = append(views, msgStylesMap[s.statusMsgType].Render(s.msg))
	}

	return statusBarStyle.Width(s.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, views...),
	)
}

func (s *StatusView) HandleWindowSize(msg tea.WindowSizeMsg) {
	s.width = msg.Width
}
