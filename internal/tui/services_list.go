package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	ServicesListKeyMap struct {
		Enter key.Binding
	}
	ServicesListItem struct {
		Name string
	}
	ServicesListView struct {
		keyMap   ServicesListKeyMap
		commands *Commands
		view     list.Model
	}
)

func (i ServicesListItem) FilterValue() string {
	return i.Name
}

func (i ServicesListItem) Title() string {
	return i.Name
}

func (i ServicesListItem) Description() string {
	return ""
}

func NewServicesListView(commands *Commands) *ServicesListView {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	view := list.New([]list.Item{}, delegate, 0, 0)
	view.Title = "Services"

	return &ServicesListView{
		keyMap:   ServicesListKeyMap{Enter: key.NewBinding(key.WithKeys("enter"))},
		commands: commands,
		view:     view,
	}
}

func (m *ServicesListView) Init() tea.Cmd {
	return nil
}

func (m *ServicesListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Enter) {
			return m, m.commands.LoadMethods(m.view.SelectedItem().(ServicesListItem).Name)
		}
	case ShowServicesList:
		items := []list.Item{}
		for _, service := range msg.Services {
			items = append(items, ServicesListItem{Name: service})
		}
		cmds = append(cmds, m.view.SetItems(items))
		cmds = append(cmds, m.commands.SetStatusOK())
	}
	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ServicesListView) View() string {
	return m.view.View()
}

func (m *ServicesListView) HandleWindowSize(msg tea.WindowSizeMsg) {
	m.view.SetWidth(msg.Width)
	m.view.SetHeight(msg.Height)
}
