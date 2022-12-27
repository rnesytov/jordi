package tui

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
	SetStatus struct {
		Status string
	}
	SetStatusMsg struct {
		Type StatusMsgType
		Msg  string
	}
	ClearStatusMsg struct{}
	ShowRequester  struct {
		Method        string
		InDescription string
		InExample     string
	}
	ShowResponse struct {
		Response string
	}
)
