package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	MethodsListKeyMap struct {
		Enter key.Binding
	}
	MethodsListItem struct {
		Name      string
		ShortName string
	}
	MethodsListView struct {
		keyMap   MethodsListKeyMap
		commands *Commands
		view     list.Model
	}
)

func NewMethodsListItem(name string) MethodsListItem {
	return MethodsListItem{
		Name:      name,
		ShortName: getShortMethodName(name),
	}
}

func (i MethodsListItem) FilterValue() string {
	return i.ShortName
}

func (i MethodsListItem) Title() string {
	return i.ShortName
}

func (i MethodsListItem) Description() string {
	return ""
}

func NewMethodsListView(commands *Commands) *MethodsListView {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	view := list.New([]list.Item{}, delegate, 0, 0)
	view.Title = "Methods"

	return &MethodsListView{
		keyMap:   MethodsListKeyMap{Enter: key.NewBinding(key.WithKeys("enter"))},
		commands: commands,
		view:     view,
	}
}

func (m *MethodsListView) Init() tea.Cmd {
	return nil
}

func (m *MethodsListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Enter) {
			return m, m.commands.LoadMethodMetadata(m.view.SelectedItem().(MethodsListItem).Name)
		}
	case ShowMethodsList:
		m.view.Title = fmt.Sprintf("Methods of %s", msg.Service)
		items := []list.Item{}
		for _, methods := range msg.Methods {
			items = append(items, NewMethodsListItem(methods))
		}
		cmds = append(cmds, m.view.SetItems(items))
		cmds = append(cmds, m.commands.SetStatusOK())
	}
	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MethodsListView) View() string {
	return m.view.View()
}

func (m *MethodsListView) HandleWindowSize(msg tea.WindowSizeMsg) {
	m.view.SetWidth(msg.Width)
	m.view.SetHeight(msg.Height)
}
