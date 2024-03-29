package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/profx5/jordi/internal/config"
	"github.com/profx5/jordi/internal/grpc"
	"github.com/profx5/jordi/internal/store"
)

type View int

const (
	Services View = iota
	Methods  View = iota
	Request  View = iota
	Response View = iota

	statusBarHeight = 1
)

type (
	RootKeyMap struct {
		Back      key.Binding
		ForceQuit key.Binding
	}
	Root struct {
		initMethod       string
		keyMap           RootKeyMap
		commands         *Commands
		currentView      View
		servicesListView *ServicesListView
		methodsListView  *MethodsListView
		requestView      *RequestView
		responseView     *ResponseView
		statusView       *StatusView
	}
)

func NewRoot(config config.Config, grpc *grpc.Wrapper, store *store.Store) *Root {
	commands := NewCommands(grpc, store)
	return &Root{
		initMethod: config.Method,
		keyMap: RootKeyMap{
			Back:      key.NewBinding(key.WithKeys("esc")),
			ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
		},
		commands:         commands,
		currentView:      Services,
		servicesListView: NewServicesListView(commands),
		methodsListView:  NewMethodsListView(commands),
		requestView:      NewRequesterView(commands),
		responseView:     NewResponseView(commands),
		statusView:       NewStatusView(),
	}
}

func (m *Root) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	if m.initMethod != "" {
		m.currentView = Request
		cmds = append(cmds, m.commands.LoadMethodMetadata(m.initMethod))
	} else {
		cmds = append(cmds, m.commands.LoadServices())
	}
	return tea.Batch(cmds...)
}

func (m *Root) CurrentView() tea.Model {
	switch m.currentView {
	case Services:
		return m.servicesListView
	case Methods:
		return m.methodsListView
	case Request:
		return m.requestView
	case Response:
		return m.responseView
	}
	panic("Unknown view")
}

func (m *Root) UpdateCurrentView(msg tea.Msg) tea.Cmd {
	switch m.currentView {
	case Services:
		updModel, cmd := m.servicesListView.Update(msg)
		m.servicesListView = updModel.(*ServicesListView)
		return cmd
	case Methods:
		updModel, cmd := m.methodsListView.Update(msg)
		m.methodsListView = updModel.(*MethodsListView)
		return cmd
	case Request:
		updModel, cmd := m.requestView.Update(msg)
		m.requestView = updModel.(*RequestView)
		return cmd
	case Response:
		updModel, cmd := m.responseView.Update(msg)
		m.responseView = updModel.(*ResponseView)
		return cmd
	}
	return nil
}

func (m *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case NewStatus, NewStatusMessage, ClearStatusMsg:
		m.statusView.Update(msg)
	case ShowServicesList:
		m.currentView = Services
	case ShowMethodsList:
		m.currentView = Methods
	case ShowRequester:
		m.currentView = Request
	case ShowResponseView:
		m.currentView = Response
	case ResendRequest:
		m.currentView = Request
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.ForceQuit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keyMap.Back) {
			cmds = append(cmds, m.UpdateCurrentView(Back{}))
			switch m.currentView {
			case Services:
				return m, tea.Quit
			case Methods:
				m.currentView = Services
			case Request:
				// exit if we have been called with a method
				if m.initMethod != "" {
					return m, tea.Quit
				}
				m.currentView = Methods
			case Response:
				m.currentView = Request
			}
			return m, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		msg.Height -= statusBarHeight
		m.servicesListView.HandleWindowSize(msg)
		m.methodsListView.HandleWindowSize(msg)
		m.requestView.HandleWindowSize(msg)
		m.responseView.HandleWindowSize(msg)
		m.statusView.HandleWindowSize(msg)
	case Err:
		cmds = append(cmds, m.commands.SetStatusMessage(msg.Error.Error(), StatusMsgError))
	}
	cmds = append(cmds, m.UpdateCurrentView(msg))
	return m, tea.Batch(cmds...)
}

func (m *Root) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.CurrentView().View(), m.statusView.View())
}
