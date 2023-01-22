package tui

import tea "github.com/charmbracelet/bubbletea"

type (
	Back             struct{}
	ShowServicesList struct {
		Services []string
	}
	ChosenService struct {
		Service string
	}
	ShowMethodsList struct {
		Service string
		Methods []string
	}
	ChosenMethod struct {
		Method string
	}
	Err struct {
		Error error
	}
	NewStatus struct {
		Status string
		Type   StatusType
	}
	NewStatusMessage struct {
		Type StatusMsgType
		Msg  string
	}
	ClearStatusMsg struct{}
	ShowRequester  struct {
		Method        string
		InDescription string
		InExample     string
	}
	ShowResponseView struct {
		ch <-chan tea.Msg
	}
	ReceivedResponse struct {
		ch       <-chan tea.Msg
		Response string
	}
	ReceivedStatus struct {
		ch     <-chan tea.Msg
		Status string
	}
	ResendRequest struct {
	}
)
