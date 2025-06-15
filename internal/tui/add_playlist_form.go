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

func updateAddPlaylistScreen(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			// Return to MenuModel
			return NewMenuModel(m.cfg), nil
		case "tab", "shift+tab":
			m.addForm.focus = 1 - m.addForm.focus
		case "up", "k":
			m.addForm.focus = 0
		case "down", "j":
			m.addForm.focus = 1
		case "enter":
			if m.addForm.name == "" || m.addForm.path == "" {
				m.addForm.errMsg = "Both fields are required"
				return m, nil
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
			// Return to MenuModel after adding playlist
			return NewMenuModel(m.cfg), nil
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
	return m, nil
}

func viewAddPlaylistScreen(m model) string {
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
