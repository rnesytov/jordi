package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusMsgType int
type StatusType int

const (
	StatusMsgError   StatusMsgType = iota
	StatusMsgSuccess StatusMsgType = iota

	StatusTypeOK    StatusType = iota
	StatusTypeWarn  StatusType = iota
	StatusTypeError StatusType = iota

	statusMsgAddWidth = 8
)

var (
	statusBackgorundColor = lipgloss.Color("#555753")
	statusBarStyle        = lipgloss.NewStyle().
				Background(statusBackgorundColor)
	statusStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true)
	statusTypeColorMap = map[StatusType]lipgloss.Color{
		StatusTypeOK:    lipgloss.Color("#4e9a06"),
		StatusTypeWarn:  lipgloss.Color("#e69b00"),
		StatusTypeError: lipgloss.Color("#ff0000"),
	}
	msgStylesMap = map[StatusMsgType]lipgloss.Style{
		StatusMsgSuccess: lipgloss.NewStyle().
			Background(statusBackgorundColor).
			Foreground(lipgloss.Color("#16a402")).
			Padding(0, 2),
		StatusMsgError: lipgloss.NewStyle().
			Background(statusBackgorundColor).
			Foreground(lipgloss.Color("#ff0000")).
			Padding(0, 2),
	}
)

type (
	StatusView struct {
		currentStatus string
		statusType    StatusType
		msg           string
		statusMsgType StatusMsgType

		width int
	}
)

func NewStatusView() *StatusView {
	return &StatusView{
		currentStatus: "Ready",
		statusType:    StatusTypeOK,
		msg:           "",
		statusMsgType: StatusMsgSuccess,
		width:         0,
	}
}

func (s *StatusView) Init() tea.Cmd {
	return nil
}

func (s *StatusView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NewStatus:
		s.currentStatus = msg.Status
		s.statusType = msg.Type
	case NewStatusMessage:
		s.msg = msg.Msg
		s.statusMsgType = msg.Type
	case ClearStatusMsg:
		s.msg = ""
	}
	return s, nil
}

func (s *StatusView) View() string {
	statusColor, ok := statusTypeColorMap[s.statusType]
	if !ok {
		statusColor = lipgloss.Color("#ffffff")
	}
	views := []string{statusStyle.Background(statusColor).Render(s.currentStatus)}
	if s.msg != "" {
		msg := s.msg
		maxWidth := s.width - len(s.currentStatus) - statusMsgAddWidth
		if maxWidth > 0 && len(msg) > maxWidth {
			msg = msg[:maxWidth]
		}
		views = append(views, msgStylesMap[s.statusMsgType].Render(msg))
	}

	return statusBarStyle.Width(s.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, views...),
	)
}

func (s *StatusView) HandleWindowSize(msg tea.WindowSizeMsg) {
	s.width = msg.Width
}
