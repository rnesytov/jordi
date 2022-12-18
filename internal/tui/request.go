package tui

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	RequestKeyMap struct {
		Send   key.Binding
		Format key.Binding
	}
	RequestView struct {
		keyMap    RequestKeyMap
		commands  *Commands
		inputView textarea.Model

		method  string
		inDesc  string
		outDesc string
	}
)

func NewRequesterView(commands *Commands) *RequestView {
	return &RequestView{
		keyMap: RequestKeyMap{
			Send:   key.NewBinding(key.WithKeys(`ctrl+s`)),
			Format: key.NewBinding(key.WithKeys(`ctrl+f`)),
		},
		commands:  commands,
		inputView: textarea.New(),
	}
}

func (r *RequestView) Init() tea.Cmd {
	return nil
}

func (r *RequestView) FormatInput() {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(r.inputView.Value()), &data)
	if err != nil {
		return
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	r.inputView.SetValue(string(b))
}

func (r *RequestView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, r.keyMap.Send) {
			return r, r.commands.SendRequest(r.method, r.inputView.Value())
		}
		if key.Matches(msg, r.keyMap.Format) {
			r.FormatInput()
		}
	case ShowRequester:
		r.method = msg.Method
		r.inDesc = msg.InDescription
		r.outDesc = msg.OutDescription

		r.inputView.Reset()
		r.inputView.SetValue(`{}`)
		r.inputView.SetCursor(1)
		r.inputView.Focus()
	}
	updInput, cmd := r.inputView.Update(msg)
	r.inputView = updInput
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
}

func (r *RequestView) View() string {
	view := strings.Builder{}
	view.WriteString(r.inputView.View())

	return view.String()
}

func (m *RequestView) HandleWindowSize(msg tea.WindowSizeMsg) {
	m.inputView.SetWidth(msg.Width)
	m.inputView.SetHeight(msg.Height - 30)
}
