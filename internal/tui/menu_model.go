package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/playlist"
)

// MenuModel is an isolated tea.Model for the menu screen.
type MenuModel struct {
	cfg    *config.Config
	menu   []menuItem
	cursor int
	quit   bool
}

func NewMenuModel(cfg *config.Config) *MenuModel {
	menu := []menuItem{
		{label: "Add playlist"},
	}
	for _, pl := range cfg.Playlists {
		menu = append(menu, menuItem{label: pl.Name})
	}
	return &MenuModel{
		cfg:    cfg,
		menu:   menu,
		cursor: 0,
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.menu)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.cursor == 0 {
				// Add playlist
				return NewAddPlaylistModel(m.cfg), nil
			}
			// Select playlist (index in cfg.Playlists is m.cursor-1)
			plIndex := m.cursor - 1
			if plIndex >= 0 && plIndex < len(m.cfg.Playlists) {
				pl := m.cfg.Playlists[plIndex]
				chans, err := playlist.LoadM3U(pl.Path)
				return NewChannelListModel(m.cfg, plIndex, pl.Name, chans, err), nil
			}
		}
	}
	// Always rebuild the menu from config to avoid duplicate entries
	m.menu = []menuItem{{label: "Add playlist"}}
	for _, pl := range m.cfg.Playlists {
		m.menu = append(m.menu, menuItem{label: pl.Name})
	}
	return m, nil
}

func (m *MenuModel) View() string {
	if m.quit {
		return "Goodbye!\n"
	}
	var b string
	b += "Vibecast\n\n"
	for i, item := range m.menu {
		cursor := "  "
		if m.cursor == i {
			cursor = "\033[7m➜\033[0m "
		}
		b += fmt.Sprintf("%s%s\n", cursor, item.label)
	}
	b += "\n[j/k] move  [enter] select  [ctrl+c] quit"
	return b
}
