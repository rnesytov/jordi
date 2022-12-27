package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/profx5/jordi/internal/grpc"
)

type Commands struct {
	grpc *grpc.GRPCWrapper
}

func NewCommands(grpc *grpc.GRPCWrapper) *Commands {
	return &Commands{grpc: grpc}
}

func (c *Commands) LoadServices() tea.Cmd {
	return func() tea.Msg {
		servicesList, err := c.grpc.ListServices()
		if err != nil {
			return Err{Error: err}
		}
		return ShowServicesList{Services: servicesList}
	}
}

func (c *Commands) LoadMethods(service string) tea.Cmd {
	return func() tea.Msg {
		methodsList, err := c.grpc.ListMethods(service)
		if err != nil {
			return Err{Error: err}
		}
		return ShowMethodsList{Service: service, Methods: methodsList}
	}
}

func (c *Commands) ShowRequester(method string) tea.Cmd {
	return func() tea.Msg {
		in, inExample, err := c.grpc.GetInOutDescription(method)
		if err != nil {
			return Err{Error: err}
		}
		return ShowRequester{
			Method:        method,
			InDescription: in,
			InExample:     inExample,
		}
	}
}

func (c *Commands) SendRequest(method string, in string) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.grpc.Invoke(method, in)
		if err != nil {
			return Err{Error: err}
		}
		return ShowResponse{Response: resp}
	}
}

func (c *Commands) SetStatus(status string) tea.Cmd {
	return func() tea.Msg {
		return SetStatus{Status: status}
	}
}
