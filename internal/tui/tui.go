package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/favourites"
	"github.com/pedrobarco/vibecast/internal/playlist"
)

type menuItem struct {
	label string
}

type mode int

const (
	modeAddPlaylist mode = iota
	modeChannelList
	modeChannelSearchInput
)

type model struct {
	cfg      *config.Config
	menu     []menuItem
	cursor   int
	quitting bool
	mode     mode
	addForm  addPlaylistForm

	// Channel list mode
	channels  []playlist.Channel
	chCursor  int
	chPlIndex int    // index of selected playlist in cfg.Playlists
	chPlName  string // name of selected playlist
	chErr     string // error loading channels

	// Filtering/bookmark modifiers
	searchQuery       string
	showOnlyBookmarks bool
	filtered          []playlist.Channel
	searchCursor      int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
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

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	switch m.mode {
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
