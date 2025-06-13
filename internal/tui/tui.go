package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
)

type menuItem struct {
	label string
}

type model struct {
	cfg        *config.Config
	menu       []menuItem
	cursor     int
	quitting   bool
}

func Run(cfg *config.Config) (tea.Model, error) {
	return tea.NewProgram(initialModel(cfg)).Run()
}

func initialModel(cfg *config.Config) model {
	menu := []menuItem{
		{label: "Add playlist"},
	}
	for _, pl := range cfg.Playlists {
		menu = append(menu, menuItem{label: pl.Name})
	}
	return model{
		cfg:    cfg,
		menu:   menu,
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.menu)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	var b strings.Builder
	b.WriteString("Vibecast\n\n")
	for i, item := range m.menu {
		cursor := "  "
		if m.cursor == i {
			cursor = "➜ "
		}
		fmt.Fprintf(&b, "%s%s\n", cursor, item.label)
	}
	b.WriteString("\n[j/k] move  [ctrl+c] quit")
	return b.String()
}
