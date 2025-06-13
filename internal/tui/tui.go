package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
)

func Run(cfg *config.Config) error {
	// Placeholder: will implement TUI logic here
	return bubbletea.NewProgram(initialModel(cfg)).Start()
}

type model struct {
	cfg *config.Config
}

func initialModel(cfg *config.Config) model {
	return model{cfg: cfg}
}

func (m model) Init() bubbletea.Cmd { return nil }
func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) { return m, nil }
func (m model) View() string { return "Vibecast TUI coming soon..." }
