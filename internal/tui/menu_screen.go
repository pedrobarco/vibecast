package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/playlist"
)

func updateMenuScreen(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			if m.cursor == 0 {
				// Add playlist
				m.mode = modeAddPlaylist
				m.addForm = addPlaylistForm{}
				return m, nil
			}
			// Select playlist (index in cfg.Playlists is m.cursor-1)
			plIndex := m.cursor - 1
			if plIndex >= 0 && plIndex < len(m.cfg.Playlists) {
				pl := m.cfg.Playlists[plIndex]
				chans, err := playlist.LoadM3U(pl.Path)
				m.channels = chans
				m.chCursor = 0
				m.chPlIndex = plIndex
				m.chPlName = pl.Name
				m.chErr = ""
				if err != nil {
					m.chErr = fmt.Sprintf("Failed to load playlist: %v", err)
				}
				// Reset menu cursor to the selected playlist for when we return
				m.cursor = m.cursor
				m.mode = modeChannelList
				return m, nil
			}
		}
	}
	// Always rebuild the menu from config to avoid duplicate entries
	// Only do this if not in channel list mode
	if m.mode != modeChannelList {
		m.menu = []menuItem{{label: "Add playlist"}}
		for _, pl := range m.cfg.Playlists {
			m.menu = append(m.menu, menuItem{label: pl.Name})
		}
	}
	return m, nil
}

func viewMenuScreen(m model) string {
	var b strings.Builder
	b.WriteString("Vibecast\n\n")
	for i, item := range m.menu {
		cursor := "  "
		if m.cursor == i {
			// Use a visible cursor for the selected item
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s\n", cursor, item.label)
	}
	b.WriteString("\n[j/k] move  [enter] select  [ctrl+c] quit")
	return b.String()
}
