package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

type keybind struct {
	Key         string
	Description string
}

type Element interface {
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	fmt.Stringer
}

type page struct {
	prev      *page
	header    string
	body      []Element
	footer    []*keybind
	cursor    int
	quit      bool
	paginator paginator.Model
}

var _ tea.Model = (*page)(nil)

func (p *page) MoveCursorUp() {
	if p.cursor > 0 {
		p.cursor--
	}
}

func (p *page) MoveCursorDown() {
	n := p.paginator.ItemsOnPage(len(p.body))
	if p.cursor < n-1 {
		p.cursor++
	}
}

func (p *page) PrevPage() {
	prev := p.paginator.Page
	p.paginator.PrevPage()
	if prev != p.paginator.Page {
		p.cursor = 0
	}
}

func (p *page) NextPage() {
	prev := p.paginator.Page
	p.paginator.NextPage()
	if prev != p.paginator.Page {
		p.cursor = 0
	}
}

func (p page) String() string {
	var b strings.Builder

	// Header
	b.WriteString(p.header + "\n\n")

	// Body
	start, end := p.paginator.GetSliceBounds(len(p.body))
	for i, item := range p.body[start:end] {
		if i == p.cursor {
			b.WriteString("> " + item.String() + "\n")
		} else {
			b.WriteString("  " + item.String() + "\n")
		}
	}

	if len(p.body) > p.paginator.PerPage {
		b.WriteString("  " + p.paginator.View())
	}

	// Footer
	b.WriteString("\n")
	for _, kb := range p.footer {
		b.WriteString("[" + kb.Key + "] " + kb.Description + "  ")
	}
	b.WriteString("\n")
	return b.String()
}

func (p *page) Init() tea.Cmd {
	return nil
}

func (p *page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := p.body[p.cursor].Update(msg)
	if m != nil {
		return m, cmd
	}
	if cmd != nil {
		return p, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			p.MoveCursorUp()
		case "j":
			p.MoveCursorDown()
		case "h":
			p.PrevPage()
		case "l":
			p.NextPage()
		case "ctrl+c":
			p.quit = true
			return p, tea.Quit
		case "esc":
			if p.prev != nil {
				return p.prev, nil
			}
		}
	}

	return p, cmd
}

func (p *page) View() string {
	if p.quit {
		return "Goodbye!\n"
	}
	return p.String()
}
