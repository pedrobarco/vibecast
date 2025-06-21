package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/favourites"
	"github.com/pedrobarco/vibecast/internal/playlist"
)

var (
	keybindM = &keybind{
		Key:         "m",
		Description: "mark/unmark",
	}
	keybindB = &keybind{
		Key:         "b",
		Description: "bookmarks",
	}
	keybindAI = &keybind{
		Key:         "a/i",
		Description: "focus",
	}
	keybindD = &keybind{
		Key:         "d",
		Description: "clear",
	}
	keybindJK = &keybind{
		Key:         "j/k",
		Description: "move",
	}
	keybindEnter = &keybind{
		Key:         "enter",
		Description: "select",
	}
	keybindCtrlC = &keybind{
		Key:         "ctrl+c",
		Description: "quit",
	}
	keybindEsc = &keybind{
		Key:         "esc",
		Description: "cancel",
	}
)

func Run(cfg *config.Config) (tea.Model, error) {
	return tea.NewProgram(NewMainModel(cfg)).Run()
}

func NewPaginator(total int) paginator.Model {
	p := paginator.New(
		paginator.WithPerPage(10),
		paginator.WithTotalPages(total),
	)
	p.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("h"),
	)
	p.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("l"),
	)
	return p
}

type mainPage struct {
	cfg *config.Config
	*page
}

func NewMainModel(cfg *config.Config) *mainPage {
	model := &mainPage{
		cfg: cfg,
		page: &page{
			prev:   nil,
			header: "Vibecast",
			body:   []Element{},
			footer: []*keybind{
				keybindJK,
				keybindEnter,
				keybindCtrlC,
			},
			cursor:    0,
			paginator: NewPaginator(0),
			quit:      false,
		},
	}

	addPlaylist := &button{
		Label: "Add Playlist",
		OnClick: func() (tea.Model, tea.Cmd) {
			next := NewAddPlaylistModel(cfg)
			next.prev = model.page
			return next, nil
		},
	}
	model.body = append(model.body, addPlaylist)

	for _, pl := range cfg.Playlists {
		item := &button{
			Label: pl.Name,
			OnClick: func() (tea.Model, tea.Cmd) {
				next, err := NewChannelListPage(cfg, cfg.Playlists[model.cursor-1])
				if err != nil {
					// TODO: create error notification cmd
					return nil, tea.Batch(
						tea.Printf("Error creating channel list page for playlist %s: %v", pl.Name, err),
					)
				}
				return next, nil
			},
		}
		model.body = append(model.body, item)
	}
	return model
}

// add playlist model
type addPlaylistPage struct {
	cfg *config.Config
	*page
}

func NewAddPlaylistModel(cfg *config.Config) *addPlaylistPage {
	model := &addPlaylistPage{
		cfg: cfg,
		page: &page{
			header: "Vibecast - Add Playlist",
			body:   []Element{},
			footer: []*keybind{
				keybindAI,
				keybindD,
				keybindJK,
				keybindEnter,
				keybindEsc,
				keybindCtrlC,
			},
			cursor:    0,
			paginator: NewPaginator(0),
		},
	}

	nameInput := NewInput("Name")
	model.body = append(model.body, nameInput)
	urlInput := NewInput("File/URL")
	model.body = append(model.body, urlInput)

	submitButton := &button{
		Label: "Submit",
		OnClick: func() (tea.Model, tea.Cmd) {
			name := nameInput.model.Value()
			url := urlInput.model.Value()
			if name == "" || url == "" {
				return nil, tea.Batch(
					tea.Println("Both fields are required"),
				)
			}

			p, err := playlist.LoadM3U(url)
			if err != nil {
				// TODO: create error notification cmd
				return nil, tea.Batch(
					tea.Printf("Error loading playlist from %s: %v", url, err),
				)
			}

			cfg.Playlists = append(cfg.Playlists, config.Playlist{
				Name: name,
				Path: url,
			})
			if err := config.Save(cfg); err != nil {
				// TODO: create error notification cmd
				return nil, tea.Batch(
					tea.Printf("Error saving config: %v", err),
				)
			}

			// TODO: create success notification cmd
			return NewMainModel(cfg), tea.Batch(
				tea.Printf("Added playlist %s with %d channels", name, len(p)),
			)
		},
	}
	model.body = append(model.body, submitButton)
	return model
}

type channelListPage struct {
	cfg      *config.Config
	pl       *config.Playlist
	channels []playlist.Channel
	*page
}

func NewChannelListPage(cfg *config.Config, pl config.Playlist) (*channelListPage, error) {
	channels, err := playlist.LoadM3U(pl.Path)
	if err != nil {
		return nil, err
	}

	model := &channelListPage{
		cfg:      cfg,
		pl:       &pl,
		channels: channels,
		page: &page{
			header: fmt.Sprintf("Vibecast - Channel List - %s (%d channels)", pl.Name, len(channels)),
			body:   []Element{},
			footer: []*keybind{
				keybindJK,
				keybindEnter,
				keybindM,
				keybindB,
				keybindEsc,
				keybindCtrlC,
			},
			cursor: 0,
		},
	}

	for _, ch := range channels {
		b := &channelButton{config: cfg, channel: &ch, playlist: &pl}
		model.body = append(model.body, b)
	}
	model.paginator = NewPaginator(len(model.body))
	return model, nil
}

func (cl *channelListPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			next := NewBookmarksPage(cl.cfg, cl.pl, cl.channels)
			next.prev = cl.page
			return next, nil
		case "/":
		}
	}

	next, cmd := cl.page.Update(msg)
	if next != cl.page {
		return next, cmd
	}
	return cl, cmd
}

func (cl *channelListPage) String() string {
	return cl.page.String()
}

type bookmarksPage struct {
	*page
}

func NewBookmarksPage(cfg *config.Config, pl *config.Playlist, channels []playlist.Channel) *bookmarksPage {
	model := &bookmarksPage{
		page: &page{
			header: fmt.Sprintf("Vibecast - Bookmarks - %s", pl.Name),
			body:   []Element{},
			footer: []*keybind{
				keybindJK,
				keybindEnter,
				keybindM,
				keybindEsc,
				keybindCtrlC,
			},
			cursor: 0,
		},
	}

	for _, ch := range channels {
		if favourites.IsFavourite(cfg.Favourites, pl.Name, ch.Name) {
			b := &channelButton{config: cfg, channel: &ch, playlist: pl}
			model.body = append(model.body, b)
		}
	}
	model.paginator = NewPaginator(len(model.body))

	return model
}
