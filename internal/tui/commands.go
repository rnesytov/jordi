package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/profx5/jordi/internal/grpc"
)

type Commands struct {
	cancel chan struct{}
	grpc   *grpc.Wrapper
}

func NewCommands(grpc *grpc.Wrapper) *Commands {
	return &Commands{grpc: grpc, cancel: make(chan struct{})}
}

func (c *Commands) LoadServices() tea.Cmd {
	return func() tea.Msg {
		select {
		case <-c.cancel:
			return nil
		case r := <-c.grpc.ListServices():
			if r.Err != nil {
				return Err{Error: r.Err}
			}
			return ShowServicesList{Services: r.Result}
		}
	}
}

func (c *Commands) LoadMethods(service string) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-c.cancel:
			return nil
		case r := <-c.grpc.ListMethods(service):
			if r.Err != nil {
				return Err{Error: r.Err}
			}
			return ShowMethodsList{Service: "", Methods: r.Result}
		}
	}
}

func (c *Commands) ShowRequester(method string) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-c.cancel:
			return nil
		case description := <-c.grpc.GetInputDescription(method):
			if description.Err != nil {
				return Err{Error: description.Err}
			}
			return ShowRequester{
				Method:        method,
				InDescription: description.Desc,
				InExample:     description.Example,
			}
		}
	}
}

func mapRespChanToMsg(ch <-chan grpc.Event) <-chan tea.Msg {
	out := make(chan tea.Msg)
	go func() {
		for respPart := range ch {
			if respPart.Err != nil {
				out <- Err{Error: respPart.Err}
			}
			switch respPart.Type {
			case grpc.EventError:
				out <- Err{Error: respPart.Err}
			case grpc.ResponseReceived:
				response := respPart.Payload.(string)
				out <- ReceivedResponse{Response: response, ch: out}
			case grpc.ReceivedTrailers:
				status := respPart.Payload.(string)
				out <- ReceivedStatus{Status: status, ch: out}
			}
		}
		close(out)
	}()
	return out
}

func (c *Commands) SendRequest(method string, payload string) tea.Cmd {
	return func() tea.Msg {
		err := checkJSON(payload)
		if err != nil {
			return Err{Error: err}
		}

		ch, err := c.grpc.Invoke(method, payload)
		if err != nil {
			return Err{Error: err}
		}
		return ShowResponseView{mapRespChanToMsg(ch)}
	}
}

func (c *Commands) SetStatus(status string, st StatusType) tea.Cmd {
	return func() tea.Msg {
		return SetStatus{Status: status, Type: st}
	}
}

func (c *Commands) SetStatusMessage(msg string, st StatusMsgType) tea.Cmd {
	return func() tea.Msg {
		return SetStatusMessage{Msg: msg, Type: st}
	}
}

func (c *Commands) ClearStatusMsg() tea.Cmd {
	return func() tea.Msg {
		return ClearStatusMsg{}
	}
}

func (c *Commands) Resend() tea.Cmd {
	return func() tea.Msg {
		return ResendRequest{}
	}
}
