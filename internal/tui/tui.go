package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/playlist"
	"github.com/pedrobarco/vibecast/internal/player"
)

type menuItem struct {
	label string
}

type mode int

const (
	modeMenu mode = iota
	modeAddPlaylist
	modeChannelList
)

type addPlaylistForm struct {
	name      string
	path      string
	focus     int // 0 = name, 1 = path
	errMsg    string
	submitted bool
}

type model struct {
	cfg         *config.Config
	menu        []menuItem
	cursor      int
	quitting    bool
	mode        mode
	addForm     addPlaylistForm

	// Channel list mode
	channels    []playlist.Channel
	chCursor    int
	chPlIndex   int    // index of selected playlist in cfg.Playlists
	chPlName    string // name of selected playlist
	chErr       string // error loading channels
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
		mode:   modeMenu,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modeMenu:
		return m.updateMenu(msg)
	case modeAddPlaylist:
		return m.updateAddPlaylist(msg)
	case modeChannelList:
		return m.updateChannelList(msg)
	default:
		return m, nil
	}
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	m.menu = []menuItem{{label: "Add playlist"}}
	for _, pl := range m.cfg.Playlists {
		m.menu = append(m.menu, menuItem{label: pl.Name})
	}
	return m, nil
}

func (m model) updateAddPlaylist(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
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
			return m, nil
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
			// Rebuild menu from config to avoid duplicate entries
			m.menu = []menuItem{{label: "Add playlist"}}
			for _, pl := range m.cfg.Playlists {
				m.menu = append(m.menu, menuItem{label: pl.Name})
			}
			m.mode = modeMenu
			// Reset cursor to the newly added playlist
			m.cursor = len(m.menu) - 1
			return m, nil
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

func (m model) updateChannelList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = modeMenu
			m.chCursor = 0
			m.channels = nil
			m.chPlName = ""
			m.chErr = ""
			return m, nil
		case "j", "down":
			if m.chCursor < len(m.channels)-1 {
				m.chCursor++
			}
		case "k", "up":
			if m.chCursor > 0 {
				m.chCursor--
			}
		case "enter":
			if m.chCursor >= 0 && m.chCursor < len(m.channels) && m.chErr == "" {
				_ = player.PlayWithVLC(m.channels[m.chCursor].URL)
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	switch m.mode {
	case modeMenu:
		return m.viewMenu()
	case modeAddPlaylist:
		return m.viewAddPlaylist()
	case modeChannelList:
		return m.viewChannelList()
	default:
		return ""
	}
}

func (m model) viewMenu() string {
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

func (m model) viewChannelList() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Playlist: %s\n\n", m.chPlName)
	if m.chErr != "" {
		fmt.Fprintf(&b, "[!] %s\n", m.chErr)
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}
	if len(m.channels) == 0 {
		b.WriteString("No channels found.\n")
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}
	for i, ch := range m.channels {
		cursor := "  "
		if m.chCursor == i {
			// Use a visible cursor for the selected channel
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s\n", cursor, ch.Name)
	}
	b.WriteString("\n[j/k] move  [enter] play  [esc] back  [ctrl+c] quit")
	return b.String()
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
