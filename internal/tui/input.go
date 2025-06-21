package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textInput struct {
	model textinput.Model
}

var _ Element = (*textInput)(nil)

func NewInput(label string) *textInput {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = label
	ti.Width = 120
	return &textInput{
		model: ti,
	}
}

// CmdNop is a no-op command that can be used to indicate no action is needed.
// It is useful for cases where a command is expected but no action is required.
func CmdNop() tea.Msg {
	return nil
}

func (i *textInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch ch := msg.String(); ch {
		case "enter", "i", "a":
			if !i.model.Focused() {
				if ch == "i" {
					i.model.CursorStart()
				}
				return nil, i.model.Focus()
			}
		case "d":
			if !i.model.Focused() {
				i.model.SetValue("")
				return nil, CmdNop
			}
		case "esc":
			if i.model.Focused() {
				i.model.Blur()
				return nil, CmdNop
			}
		}
	}
	model, cmd := i.model.Update(msg)
	i.model = model
	return nil, cmd
}

func (i *textInput) String() string {
	return i.model.View()
}
