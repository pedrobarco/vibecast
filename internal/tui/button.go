package tui

import tea "github.com/charmbracelet/bubbletea"

type button struct {
	Label   string
	OnClick func() (tea.Model, tea.Cmd)
}

var _ Element = (*button)(nil)

func (a button) String() string {
	return a.Label
}

func (a button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return a.OnClick()
		}
	}
	return nil, nil
}
