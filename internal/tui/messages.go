package tui

type (
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
	ShowStatus struct {
		Status string
	}
	ShowRequester struct {
		Method         string
		InDescription  string
		OutDescription string
	}
	ShowResponse struct {
		Response string
	}
)
