package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type model struct {
	pageName string
	// Used in search view
	search          string
	displayArticles map[int]string
	cursor          int
	// Article view
	shownArticle string
	viewport     viewport.Model
	ready        bool
	content      string
}

func initialModel(topic string) model {
	return model{
		pageName:        "search",
		search:          topic,
		displayArticles: map[int]string{0: "Lions", 1: "India", 2: "Submarines", 3: "Turtles", 4: "Canada", 5: "Go_(programming_language)"}, //make(map[int]string)
		content:         "blahblah\nblah\nblah\nblah\nblah\nb\nl\na\nh\nblah\nblah\n...\nblah",
	}
}

func (m model) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func ArticleView(m model) string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func ArticleUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "left", "h":
			m.pageName = "search"
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func SearchView(m model) string {
	s := "wki - Search Wikipedia\n\n"
	for i := 0; i < len(m.displayArticles); i++ {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		// Render the row
		s += fmt.Sprintf("%s [%s] \n", cursor, m.displayArticles[i])
	}

	// The footer
	s += "\nPress esc to quit.\n"

	// Send the UI for rendering
	return s
}

func SearchUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.displayArticles)-1 {
				m.cursor++
			}

		case "enter", "right", "l":
			m.pageName = "article"
			m.shownArticle = m.displayArticles[m.cursor]
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

type Page struct {
	update func(model, tea.Msg) (tea.Model, tea.Cmd)
	view   func(model) string
}

var pages = map[string]Page{
	"search":  {update: SearchUpdate, view: SearchView},
	"article": {update: ArticleUpdate, view: ArticleView},
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if page, ok := pages[m.pageName]; ok {
		return page.update(m, msg)
	}
	return m, tea.Quit
}

func (m model) View() string {
	if page, ok := pages[m.pageName]; ok {
		return page.view(m)
	}
	return "I don't know how you ended up here.."
}

func main() {
	topic := flag.String("topic", "", "Topic to search")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *help {
		flag.Usage()
	}

	p := tea.NewProgram(initialModel(*topic))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
