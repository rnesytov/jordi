package tui

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	descriptionStyle = lipgloss.NewStyle().PaddingLeft(2).Border(lipgloss.NormalBorder(), true, false)
)

type (
	RequestKeyMap struct {
		Send       key.Binding
		Format     key.Binding
		ToggleDesc key.Binding
	}
	RequestView struct {
		keyMap      RequestKeyMap
		commands    *Commands
		inputView   textarea.Model
		requestDesc string
		title       TitleView
		help        HelpView

		method string
		inDesc string

		width, height int
		showDesc      bool
	}
)

func (r RequestKeyMap) Bindings() []key.Binding {
	return []key.Binding{
		r.Send,
		r.Format,
		r.ToggleDesc,
	}
}

func DefaultRequestKeyMap() RequestKeyMap {
	send := key.NewBinding(key.WithKeys("ctrl+s"))
	send.SetHelp(`ctrl+s`, "send")

	format := key.NewBinding(key.WithKeys("ctrl+f"))
	format.SetHelp(`ctrl+f`, "format")

	toggleDesc := key.NewBinding(key.WithKeys("tab"))
	toggleDesc.SetHelp(`tab`, "description")

	return RequestKeyMap{
		Send:       send,
		Format:     format,
		ToggleDesc: toggleDesc,
	}
}

func NewRequesterView(commands *Commands) *RequestView {
	inputView := textarea.New()
	inputView.ShowLineNumbers = false
	inputView.CharLimit = 0
	inputView.Prompt = ""

	keyMap := DefaultRequestKeyMap()

	return &RequestView{
		keyMap:      keyMap,
		commands:    commands,
		inputView:   inputView,
		requestDesc: "",
		title:       NewTitleView("Request"),
		help:        NewHelpView(keyMap),
		method:      "",
		inDesc:      "",
		showDesc:    false,
		height:      0,
		width:       0,
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
		} else if key.Matches(msg, r.keyMap.Format) {
			r.FormatInput()
		} else if key.Matches(msg, r.keyMap.ToggleDesc) && r.inDesc != "" {
			r.showDesc = !r.showDesc
		} else {
			cmds = append(cmds, r.commands.ClearStatusMsg())
		}
	case ShowRequester:
		r.method = msg.Method
		r.inDesc = msg.InDescription

		r.inputView.Reset()
		r.inputView.SetValue(msg.InExample)
		r.inputView.SetCursor(1)
		r.inputView.Focus()

		r.title.SetTitle(getShortMethodName(msg.Method))
		cmds = append(cmds, r.commands.SetStatusOK())
	case ResendRequest:
		return r, r.commands.SendRequest(r.method, r.inputView.Value())
	}

	updInput, cmd := r.inputView.Update(msg)
	r.inputView = updInput
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
}

func (r *RequestView) View() string {
	r.SyncSize()

	views := []string{r.title.View(), r.inputView.View()}
	if r.showDesc {
		views = append(views, descriptionStyle.Render(r.inDesc))
	}
	views = append(views, r.help.View())

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (r *RequestView) HandleWindowSize(msg tea.WindowSizeMsg) {
	r.width, r.height = msg.Width, msg.Height
}

func (r *RequestView) SyncSize() {
	r.inputView.SetWidth(r.width)
	r.help.SetWidth(r.width)

	height := r.height - helpHeight - titleHeight
	if r.showDesc {
		height = height - helpHeight - countLines(r.inDesc) - 2
	}
	r.inputView.SetHeight(height)
}
