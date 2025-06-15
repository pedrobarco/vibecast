package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/favourites"
	"github.com/pedrobarco/vibecast/internal/player"
)

// Channel List Screen (modeChannelList)
func updateChannelListScreen(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = modeMenu
			m.chCursor = 0
			m.channels = nil
			m.chPlName = ""
			m.chErr = ""
			m.searchQuery = ""
			m.showOnlyBookmarks = false
			return m, nil
		case "/":
			m.mode = modeChannelSearchInput
			// Don't reset searchQuery, allow editing
			return m, nil
		case "b":
			m.showOnlyBookmarks = !m.showOnlyBookmarks
			m.chCursor = 0
			return m, nil
		case "m":
			visible := m.visibleChannels()
			if m.chCursor >= 0 && m.chCursor < len(visible) {
				ch := visible[m.chCursor]
				if favourites.IsFavourite(m.cfg.Favourites, m.chPlName, ch.Name) {
					favourites.RemoveFavourite(m.cfg.Favourites, m.chPlName, ch.Name)
				} else {
					favourites.AddFavourite(m.cfg.Favourites, m.chPlName, ch.Name)
				}
				home := os.Getenv("HOME")
				cfgPath := home + "/.config/vibecast/config.yaml"
				_ = config.Save(cfgPath, m.cfg)
			}
			return m, nil
		case "j", "down":
			visible := m.visibleChannels()
			if m.chCursor < len(visible)-1 {
				m.chCursor++
			}
		case "k", "up":
			if m.chCursor > 0 {
				m.chCursor--
			}
		case "enter":
			visible := m.visibleChannels()
			if m.chCursor >= 0 && m.chCursor < len(visible) && m.chErr == "" {
				_ = player.PlayWithVLC(visible[m.chCursor].URL)
			}
		}
	}
	return m, nil
}

func viewChannelListScreen(m model) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Playlist: %s\n\n", m.chPlName)
	if m.chErr != "" {
		fmt.Fprintf(&b, "[!] %s\n", m.chErr)
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}
	visible := m.visibleChannels()
	if len(visible) == 0 {
		b.WriteString("No channels found.\n")
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}

	const windowSize = 15
	start, end := paginate(len(visible), m.chCursor, windowSize)

	for i := start; i < end; i++ {
		ch := visible[i]
		cursor := "  "
		star := " "
		if favourites.IsFavourite(m.cfg.Favourites, m.chPlName, ch.Name) {
			star = "★"
		}
		if m.chCursor == i {
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s %s\n", cursor, star, ch.Name)
	}
	b.WriteString(fmt.Sprintf("\nShowing %d-%d of %d channels", start+1, end, len(visible)))
	b.WriteString("\n[j/k] move  [enter] play  [/] search  [b] bookmarks  [m] mark/unmark  [esc] back  [ctrl+c] quit")
	if m.showOnlyBookmarks {
		b.WriteString(" [BOOKMARKS]")
	}
	if m.searchQuery != "" {
		b.WriteString(fmt.Sprintf(" [SEARCH: %s]", m.searchQuery))
	}
	return b.String()
}

// Channel Search Input Screen (modeChannelSearchInput)
func updateChannelSearchInputScreen(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// Cancel the filter/search
			m.searchQuery = ""
			m.mode = modeChannelList
			return m, nil
		case "j", "down":
			visible := m.visibleChannels()
			if m.chCursor < len(visible)-1 {
				m.chCursor++
			}
		case "k", "up":
			if m.chCursor > 0 {
				m.chCursor--
			}
		case "enter":
			// Accept the filter and switch to the list view
			m.mode = modeChannelList
			return m, nil
		}
	}
	return m, nil
}

func viewChannelSearchInputScreen(m model) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Search: /%s\n\n", m.searchQuery))
	visible := m.visibleChannels()
	if len(visible) == 0 {
		b.WriteString("No channels found.\n")
		b.WriteString("\n[esc] back  [ctrl+c] quit")
		return b.String()
	}
	const windowSize = 15
	start, end := paginate(len(visible), m.chCursor, windowSize)
	for i := start; i < end; i++ {
		ch := visible[i]
		cursor := "  "
		star := " "
		if favourites.IsFavourite(m.cfg.Favourites, m.chPlName, ch.Name) {
			star = "★"
		}
		if m.chCursor == i {
			cursor = "\033[7m➜\033[0m "
		}
		fmt.Fprintf(&b, "%s%s %s\n", cursor, star, ch.Name)
	}
	b.WriteString(fmt.Sprintf("\nShowing %d-%d of %d channels", start+1, end, len(visible)))
	b.WriteString("\n[j/k] move  [enter] play  [esc] back to list  [ctrl+c] quit")
	return b.String()
}
