package tui

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	requestHelpHeight = 1
)

var (
	helpStyle        = lipgloss.NewStyle().PaddingLeft(2)
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
		help        help.Model

		method string
		inDesc string

		width, height int
		showDesc      bool
	}
)

func (r *RequestKeyMap) Bindings() []key.Binding {
	return []key.Binding{
		r.Send,
		r.Format,
		r.ToggleDesc,
	}
}

func DefaultKeyMap() RequestKeyMap {
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
	inputView.Prompt = ""

	return &RequestView{
		keyMap:      DefaultKeyMap(),
		commands:    commands,
		inputView:   inputView,
		requestDesc: "",
		help:        help.New(),
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
		if key.Matches(msg, r.keyMap.ToggleDesc) && r.inDesc != "" {
			r.showDesc = !r.showDesc
		}
	case ShowRequester:
		r.method = msg.Method
		r.inDesc = msg.InDescription

		r.inputView.Reset()
		r.inputView.SetValue(msg.InExample)
		r.inputView.SetCursor(1)
		r.inputView.Focus()
	}
	updInput, cmd := r.inputView.Update(msg)
	r.inputView = updInput
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
}

func (r *RequestView) View() string {
	r.SyncSize()

	views := []string{r.inputView.View()}
	if r.showDesc {
		// md := fmt.Sprintf("```protobuf\n%s", r.inDesc)
		// renderer, _ := glamour.NewTermRenderer(
		// 	glamour.WithAutoStyle(),
		// 	glamour.WithWordWrap(r.width),
		// )
		// gl, err := renderer.Render(md)
		// if err != nil {
		// 	gl = r.inDesc
		// }
		views = append(views, descriptionStyle.Render(r.inDesc))
	}
	views = append(views, helpStyle.Render(r.help.ShortHelpView(r.keyMap.Bindings())))

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (r *RequestView) HandleWindowSize(msg tea.WindowSizeMsg) {
	r.width, r.height = msg.Width, msg.Height
}

func (r *RequestView) SyncSize() {
	r.inputView.SetHeight(r.height - requestHelpHeight)
	r.inputView.SetWidth(r.width)
	r.help.Width = r.width

	if r.showDesc {
		r.inputView.SetHeight(r.height - requestHelpHeight - CountNewLines(r.inDesc) - 3)
	}
}

func CountNewLines(s string) int {
	return strings.Count(s, "\n")
}
