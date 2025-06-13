package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
)

func Run(cfg *config.Config) (tea.Model, error) {
	// Placeholder: will implement TUI logic here
	return tea.NewProgram(initialModel(cfg)).Run()
}

type model struct {
	cfg *config.Config
}

func initialModel(cfg *config.Config) model {
	return model{cfg: cfg}
}

func (m model) Init() tea.Cmd                           { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m model) View() string                            { return "Vibecast TUI coming soon..." }
