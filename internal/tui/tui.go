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
	modeChannelSearchBrowse
	modeFavouritesList
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

	// Channel search mode
	searchQuery string
	filtered    []playlist.Channel
	searchCursor int

	// Favourites mode
	favChannels []playlist.Channel
	favCursor   int
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
	case modeChannelSearchInput:
		return m.updateChannelSearchInput(msg)
	case modeChannelSearchBrowse:
		return m.updateChannelSearchBrowse(msg)
	case modeFavouritesList:
		return m.updateFavouritesList(msg)
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
		case "/":
			m.mode = modeChannelSearchInput
			m.searchQuery = ""
			m.filtered = m.channels
			m.searchCursor = m.chCursor
			return m, nil
		case "b":
			// Show favourites for this playlist
			favs := m.cfg.Favourites[m.chPlName]
			var favChannels []playlist.Channel
			for _, ch := range m.channels {
				for _, fav := range favs {
					if ch.Name == fav {
						favChannels = append(favChannels, ch)
						break
					}
				}
			}
			m.favChannels = favChannels
			m.favCursor = 0
			m.mode = modeFavouritesList
			return m, nil
		case "m":
			// Toggle favourite for this channel
			if m.chCursor >= 0 && m.chCursor < len(m.channels) {
				ch := m.channels[m.chCursor]
				if favourites.IsFavourite(m.cfg.Favourites, m.chPlName, ch.Name) {
					favourites.RemoveFavourite(m.cfg.Favourites, m.chPlName, ch.Name)
				} else {
					favourites.AddFavourite(m.cfg.Favourites, m.chPlName, ch.Name)
				}
				// Save config
				home := os.Getenv("HOME")
				cfgPath := home + "/.config/vibecast/config.yaml"
				_ = config.Save(cfgPath, m.cfg)
			}
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
	case modeChannelSearchInput:
		return m.viewChannelSearchInput()
	case modeChannelSearchBrowse:
		return m.viewChannelSearchBrowse()
	case modeFavouritesList:
		return m.viewFavouritesList()
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

	// Pagination: show only a window of channels around the cursor
	const windowSize = 15
	start := m.chCursor - windowSize/2
	if start < 0 {
		start = 0
	}
	end := start + windowSize
	if end > len(m.channels) {
		end = len(m.channels)
	}
	if end-start < windowSize && end == len(m.channels) {
		start = end - windowSize
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		ch := m.channels[i]
		cursor := "  "
		star := " "
		if favourites.IsFavourite(m.cfg.Favourites, m.chPlName, ch.Name) {
			star = "★"
		}
		if m.chCursor == i {
			// Use a visible cursor for the selected channel
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s %s\n", cursor, star, ch.Name)
	}
	b.WriteString(fmt.Sprintf("\nShowing %d-%d of %d channels", start+1, end, len(m.channels)))
	b.WriteString("\n[j/k] move  [enter] play  [/] search  [m] mark/unmark  [b] bookmarks  [esc] back  [ctrl+c] quit")
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
func (m model) updateChannelSearchInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			m.searchQuery += msg.String()
		case tea.KeyBackspace:
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
		}
		switch msg.String() {
		case "esc":
			// If there is a query, go to browse filtered, else go back to all channels
			if m.searchQuery != "" {
				// Update filtered list before switching
				m.filtered = nil
				lq := strings.ToLower(m.searchQuery)
				for _, ch := range m.channels {
					if fuzzy.Match(lq, strings.ToLower(ch.Name)) {
						m.filtered = append(m.filtered, ch)
					}
				}
				if m.searchCursor >= len(m.filtered) {
					m.searchCursor = len(m.filtered) - 1
				}
				if m.searchCursor < 0 {
					m.searchCursor = 0
				}
				m.mode = modeChannelSearchBrowse
				return m, nil
			} else {
				m.mode = modeChannelList
				m.filtered = nil
				m.searchCursor = 0
				return m, nil
			}
		case "j", "down":
			// allow navigation while searching
			if m.searchCursor < len(m.filtered)-1 {
				m.searchCursor++
			}
		case "k", "up":
			if m.searchCursor > 0 {
				m.searchCursor--
			}
		case "enter":
			if m.searchCursor >= 0 && m.searchCursor < len(m.filtered) {
				_ = player.PlayWithVLC(m.filtered[m.searchCursor].URL)
			}
		}
		// Update filtered list
		m.filtered = nil
		lq := strings.ToLower(m.searchQuery)
		for _, ch := range m.channels {
			if fuzzy.Match(lq, strings.ToLower(ch.Name)) {
				m.filtered = append(m.filtered, ch)
			}
		}
		// Reset cursor if out of bounds
		if m.searchCursor >= len(m.filtered) {
			m.searchCursor = len(m.filtered) - 1
		}
		if m.searchCursor < 0 {
			m.searchCursor = 0
		}
	}
	return m, nil
}

func (m model) updateChannelSearchBrowse(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			m.mode = modeChannelSearchInput
			return m, nil
		case "esc":
			// Go back to all channels view
			m.mode = modeChannelList
			m.filtered = nil
			m.searchQuery = ""
			m.searchCursor = 0
			return m, nil
		case "m":
			// Toggle favourite for this channel in filtered list
			if m.searchCursor >= 0 && m.searchCursor < len(m.filtered) {
				ch := m.filtered[m.searchCursor]
				if favourites.IsFavourite(m.cfg.Favourites, m.chPlName, ch.Name) {
					favourites.RemoveFavourite(m.cfg.Favourites, m.chPlName, ch.Name)
				} else {
					favourites.AddFavourite(m.cfg.Favourites, m.chPlName, ch.Name)
				}
				// Save config
				home := os.Getenv("HOME")
				cfgPath := home + "/.config/vibecast/config.yaml"
				_ = config.Save(cfgPath, m.cfg)
			}
			return m, nil
		case "b":
			// Show favourites for this playlist
			favs := m.cfg.Favourites[m.chPlName]
			var favChannels []playlist.Channel
			for _, ch := range m.channels {
				for _, fav := range favs {
					if ch.Name == fav {
						favChannels = append(favChannels, ch)
						break
					}
				}
			}
			m.favChannels = favChannels
			m.favCursor = 0
			m.mode = modeFavouritesList
			return m, nil
		case "j", "down":
			if m.searchCursor < len(m.filtered)-1 {
				m.searchCursor++
			}
		case "k", "up":
			if m.searchCursor > 0 {
				m.searchCursor--
			}
		case "enter":
			if m.searchCursor >= 0 && m.searchCursor < len(m.filtered) {
				_ = player.PlayWithVLC(m.filtered[m.searchCursor].URL)
			}
		}
	}
	return m, nil
}

func (m model) viewChannelSearchInput() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Search: /%s\n\n", m.searchQuery)
	if len(m.filtered) == 0 {
		b.WriteString("No channels found.\n")
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}
	const windowSize = 15
	start := m.searchCursor - windowSize/2
	if start < 0 {
		start = 0
	}
	end := start + windowSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	if end-start < windowSize && end == len(m.filtered) {
		start = end - windowSize
		if start < 0 {
			start = 0
		}
	}
	for i := start; i < end; i++ {
		ch := m.filtered[i]
		cursor := "  "
		if m.searchCursor == i {
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s\n", cursor, ch.Name)
	}
	b.WriteString(fmt.Sprintf("\nShowing %d-%d of %d channels", start+1, end, len(m.filtered)))
	b.WriteString("\n[j/k] move  [enter] play  [esc] browse filtered  [ctrl+c] quit")
	return b.String()
}

func (m model) viewChannelSearchBrowse() string {
	var b strings.Builder
	b.WriteString("Filtered results (press / to search again, esc to show all):\n\n")
	if len(m.filtered) == 0 {
		b.WriteString("No channels found.\n")
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}
	const windowSize = 15
	start := m.searchCursor - windowSize/2
	if start < 0 {
		start = 0
	}
	end := start + windowSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	if end-start < windowSize && end == len(m.filtered) {
		start = end - windowSize
		if start < 0 {
			start = 0
		}
	}
	for i := start; i < end; i++ {
		ch := m.filtered[i]
		cursor := "  "
		if m.searchCursor == i {
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s\n", cursor, ch.Name)
	}
	b.WriteString(fmt.Sprintf("\nShowing %d-%d of %d channels", start+1, end, len(m.filtered)))
	b.WriteString("\n[j/k] move  [enter] play  [/] search again  [esc] show all  [ctrl+c] quit")
	return b.String()
}
func (m model) updateFavouritesList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = modeChannelList
			return m, nil
		case "j", "down":
			if m.favCursor < len(m.favChannels)-1 {
				m.favCursor++
			}
		case "k", "up":
			if m.favCursor > 0 {
				m.favCursor--
			}
		case "m":
			// Remove favourite
			if m.favCursor >= 0 && m.favCursor < len(m.favChannels) {
				ch := m.favChannels[m.favCursor]
				favourites.RemoveFavourite(m.cfg.Favourites, m.chPlName, ch.Name)
				// Save config
				home := os.Getenv("HOME")
				cfgPath := home + "/.config/vibecast/config.yaml"
				_ = config.Save(cfgPath, m.cfg)
				// Remove from view
				m.favChannels = append(m.favChannels[:m.favCursor], m.favChannels[m.favCursor+1:]...)
				if m.favCursor >= len(m.favChannels) && m.favCursor > 0 {
					m.favCursor--
				}
			}
			return m, nil
		case "enter":
			if m.favCursor >= 0 && m.favCursor < len(m.favChannels) {
				_ = player.PlayWithVLC(m.favChannels[m.favCursor].URL)
			}
		}
	}
	return m, nil
}

func (m model) viewFavouritesList() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Favourites for playlist: %s\n\n", m.chPlName)
	if len(m.favChannels) == 0 {
		b.WriteString("No favourites yet.\n")
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}
	const windowSize = 15
	start := m.favCursor - windowSize/2
	if start < 0 {
		start = 0
	}
	end := start + windowSize
	if end > len(m.favChannels) {
		end = len(m.favChannels)
	}
	if end-start < windowSize && end == len(m.favChannels) {
		start = end - windowSize
		if start < 0 {
			start = 0
		}
	}
	for i := start; i < end; i++ {
		ch := m.favChannels[i]
		cursor := "  "
		if m.favCursor == i {
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s\n", cursor, ch.Name)
	}
	b.WriteString(fmt.Sprintf("\nShowing %d-%d of %d favourites", start+1, end, len(m.favChannels)))
	b.WriteString("\n[j/k] move  [enter] play  [m] remove  [esc] back  [ctrl+c] quit")
	return b.String()
}
