package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/playlist"
	"github.com/pedrobarco/vibecast/internal/player"
	"github.com/pedrobarco/vibecast/internal/favourites"
)

type menuItem struct {
	label string
}

type mode int

const (
	modeMenu mode = iota
	modeAddPlaylist
	modeChannelList
	modeChannelSearchInput
)


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

	// Filtering/bookmark modifiers
	searchQuery string
	showOnlyBookmarks bool
	filtered    []playlist.Channel
	searchCursor int
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
		return updateAddPlaylistScreen(m, msg)
	case modeChannelList:
		return updateChannelListScreen(m, msg)
	case modeChannelSearchInput:
		return updateChannelSearchInputScreen(m, msg)
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
	// Only do this if not in channel list mode
	if m.mode != modeChannelList {
		m.menu = []menuItem{{label: "Add playlist"}}
		for _, pl := range m.cfg.Playlists {
			m.menu = append(m.menu, menuItem{label: pl.Name})
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
		return viewAddPlaylistScreen(m)
	case modeChannelList:
		return viewChannelListScreen(m)
	case modeChannelSearchInput:
		return viewChannelSearchInputScreen(m)
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

func (m model) visibleChannels() []playlist.Channel {
	chans := m.channels
	if m.showOnlyBookmarks {
		var filtered []playlist.Channel
		for _, ch := range chans {
			if favourites.IsFavourite(m.cfg.Favourites, m.chPlName, ch.Name) {
				filtered = append(filtered, ch)
			}
		}
		chans = filtered
	}
	if m.searchQuery != "" {
		var filtered []playlist.Channel
		lq := strings.ToLower(m.searchQuery)
		for _, ch := range chans {
			if fuzzy.Match(lq, strings.ToLower(ch.Name)) {
				filtered = append(filtered, ch)
			}
		}
		chans = filtered
	}
	return chans
}

func paginate(length, cursor, windowSize int) (start, end int) {
	start = cursor - windowSize/2
	if start < 0 {
		start = 0
	}
	end = start + windowSize
	if end > length {
		end = length
	}
	if end-start < windowSize && end == length {
		start = end - windowSize
		if start < 0 {
			start = 0
		}
	}
	return start, end
}





