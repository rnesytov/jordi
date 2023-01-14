package tui

import "github.com/charmbracelet/lipgloss"

const (
	titleHeight = 2
)

var (
	titleBarStyle = lipgloss.NewStyle().Padding(0, 0, 1, 2)
	titleStyle    = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)
)

type TitleView struct {
	title string
}

func NewTitleView(title string) TitleView {
	return TitleView{title: title}
}

func (v *TitleView) View() string {
	return titleBarStyle.Render(titleStyle.Render(v.title))
}

func (v *TitleView) SetTitle(title string) {
	v.title = title
}
