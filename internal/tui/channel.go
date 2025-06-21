package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/favourites"
	"github.com/pedrobarco/vibecast/internal/player"
	"github.com/pedrobarco/vibecast/internal/playlist"
)

type channelButton struct {
	config   *config.Config
	playlist *config.Playlist
	channel  *playlist.Channel
}

func (c channelButton) Favourite() bool {
	return favourites.IsFavourite(c.config.Favourites, c.playlist.Name, c.channel.Name)
}

func (c channelButton) String() string {
	var b strings.Builder
	if c.Favourite() {
		b.WriteString("* ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString(c.channel.Name)
	return b.String()
}

func (c *channelButton) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m":
			if c.Favourite() {
				favourites.RemoveFavourite(c.config.Favourites, c.playlist.Name, c.channel.Name)
			} else {
				favourites.AddFavourite(c.config.Favourites, c.playlist.Name, c.channel.Name)
			}
			if err := config.Save(c.config); err != nil {
				// TODO: create error notification cmd
				return nil, tea.Batch(
					tea.Printf("Error saving favourites: %v", err),
				)
			}
			tea.Printf("Toggled favourite for channel: %s\n", c.channel.Name)
			return nil, CmdNop
		case "enter":
			err := player.PlayWithVLC(c.channel.URL)
			if err != nil {
				// TODO: create error notification cmd
				return nil, tea.Batch(
					tea.Printf("Error playing channel %s: %v", c.channel.Name, err),
				)
			}
			return nil, CmdNop
		}
	}
	return nil, nil
}
