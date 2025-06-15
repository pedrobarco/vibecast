package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
)

type addPlaylistForm struct {
	name      string
	path      string
	focus     int // 0 = name, 1 = path
	errMsg    string
	submitted bool
}

func (m *model) updateAddPlaylist(msg any) (model, func() tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return *m, tea.Quit
		case "esc":
			m.mode = modeMenu
			// Rebuild menu from config to avoid duplicate entries
			m.menu = []menuItem{{label: "Add playlist"}}
			for _, pl := range m.cfg.Playlists {
				m.menu = append(m.menu, menuItem{label: pl.Name})
			}
			// Reset cursor to first playlist if any
			if len(m.menu) > 1 {
				m.cursor = 1
			} else {
				m.cursor = 0
			}
			return *m, nil
		case "tab", "shift+tab":
			m.addForm.focus = 1 - m.addForm.focus
		case "up", "k":
			m.addForm.focus = 0
		case "down", "j":
			m.addForm.focus = 1
		case "enter":
			if m.addForm.name == "" || m.addForm.path == "" {
				m.addForm.errMsg = "Both fields are required"
				return *m, nil
			}
			// Add playlist to config
			m.cfg.Playlists = append(m.cfg.Playlists, config.Playlist{
				Name: m.addForm.name,
				Path: m.addForm.path,
			})
			// Save config
			home := os.Getenv("HOME")
			cfgPath := home + "/.config/vibecast/config.yaml"
			_ = config.Save(cfgPath, m.cfg)
			// Rebuild menu from config to avoid duplicate entries
			m.menu = []menuItem{{label: "Add playlist"}}
			for _, pl := range m.cfg.Playlists {
				m.menu = append(m.menu, menuItem{label: pl.Name})
			}
			m.mode = modeMenu
			// Reset cursor to the newly added playlist
			m.cursor = len(m.menu) - 1
			return *m, nil
		default:
			// Handle text input
			if m.addForm.focus == 0 {
				// Name field
				if msg.Type == tea.KeyRunes {
					m.addForm.name += msg.String()
				} else if msg.Type == tea.KeyBackspace && len(m.addForm.name) > 0 {
					m.addForm.name = m.addForm.name[:len(m.addForm.name)-1]
				}
			} else {
				// Path field
				if msg.Type == tea.KeyRunes {
					m.addForm.path += msg.String()
				} else if msg.Type == tea.KeyBackspace && len(m.addForm.path) > 0 {
					m.addForm.path = m.addForm.path[:len(m.addForm.path)-1]
				}
			}
		}
	}
	return *m, nil
}

func (m model) viewAddPlaylist() string {
	var b strings.Builder
	b.WriteString("Add Playlist\n\n")
	nameLabel := "Name: "
	pathLabel := "File/URL: "
	if m.addForm.focus == 0 {
		nameLabel = "> Name: "
	} else {
		pathLabel = "> File/URL: "
	}
	fmt.Fprintf(&b, "%s%s\n", nameLabel, m.addForm.name)
	fmt.Fprintf(&b, "%s%s\n", pathLabel, m.addForm.path)
	if m.addForm.errMsg != "" {
		fmt.Fprintf(&b, "\n[!] %s\n", m.addForm.errMsg)
	}
	b.WriteString("\n[tab] switch field  [enter] submit  [esc] cancel  [ctrl+c] quit")
	return b.String()
}
