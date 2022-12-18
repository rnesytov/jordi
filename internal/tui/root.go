package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/profx5/jordi/internal/grpc"
)

type View int

const (
	Services View = iota
	Methods  View = iota
	Request  View = iota
	Response View = iota
)

type (
	RootKeyMap struct {
		Back      key.Binding
		ForceQuit key.Binding
	}
	Root struct {
		keyMap           RootKeyMap
		commands         *Commands
		currentView      View
		servicesListView *ServicesListView
		methodsListView  *MethodsListView
		requestView      *RequestView
		responseView     *ResponseView
	}
)

func NewRoot(grpc *grpc.GRPCWrapper) *Root {
	commands := NewCommands(grpc)
	return &Root{
		keyMap: RootKeyMap{
			Back:      key.NewBinding(key.WithKeys("esc")),
			ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
		},
		commands:         commands,
		currentView:      Services,
		servicesListView: NewServicesListView(commands),
		methodsListView:  NewMethodsListView(commands),
		requestView:      NewRequesterView(commands),
		responseView:     NewResponseView(),
	}
}

func (m *Root) Init() tea.Cmd {
	return m.commands.LoadServices()
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

func (m *Root) UpdateAllViews(msg tea.Cmd) tea.Cmd {
	cmds := []tea.Cmd{}

	updModel, cmd := m.servicesListView.Update(msg)
	m.servicesListView = updModel.(*ServicesListView)
	cmds = append(cmds, cmd)

	updModel, cmd = m.methodsListView.Update(msg)
	m.methodsListView = updModel.(*MethodsListView)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case ShowServicesList:
		m.currentView = Services
	case ShowMethodsList:
		m.currentView = Methods
	case ShowRequester:
		m.currentView = Request
	case ShowResponse:
		m.currentView = Response
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.ForceQuit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keyMap.Back) {
			switch m.currentView {
			case Services:
				return m, tea.Quit
			case Methods:
				m.currentView = Services
			case Request:
				m.currentView = Methods
			case Response:
				m.currentView = Request
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.servicesListView.HandleWindowSize(msg)
		m.methodsListView.HandleWindowSize(msg)
		m.requestView.HandleWindowSize(msg)
		m.responseView.HandleWindowSize(msg)
	case Err:
		panic(msg.Error)
	}
	cmds = append(cmds, m.UpdateCurrentView(msg))
	return m, tea.Batch(cmds...)
}

func (m *Root) View() string {
	doc := strings.Builder{}
	doc.WriteString(m.CurrentView().View())
	return doc.String()
}
