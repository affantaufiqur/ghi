package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ghi/repo"
)

var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#F25D94")).
		Padding(0, 1)

	issueStyle = lipgloss.NewStyle().Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)

	openColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3FB950"))
	closedColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#F85149"))
)

type Model struct {
	Issues         []repo.Issue
	FilteredIssues []repo.Issue
	Cursor         int
	Offset         int
	Repo           string
	Width          int
	Height         int
	SearchMode     bool
	Query          textinput.Model
}

func New() Model {
	q := textinput.New()
	q.Placeholder = "search issues..."
	q.Prompt = "/ "
	return Model{Query: q}
}

func (m *Model) moveCursor(delta int) {
	items := m.visibleItems()
	if len(items) == 0 {
		return
	}
	m.Cursor += delta
	if m.Cursor < 0 {
		m.Cursor = 0
	}
	if m.Cursor >= len(items) {
		m.Cursor = len(items) - 1
	}

	visible := m.Height - 3
	if m.SearchMode {
		visible--
	}
	if visible < 1 {
		visible = 1
	}

	if m.Cursor < m.Offset {
		m.Offset = m.Cursor
	}
	if m.Cursor >= m.Offset+visible {
		m.Offset = m.Cursor - visible + 1
	}
}

func (m *Model) UpdateSize(w, h int) {
	m.Width = w
	m.Height = h
	m.Query.Width = w - 4
}

func (m Model) visibleItems() []repo.Issue {
	if m.SearchMode && m.Query.Value() != "" {
		return m.FilteredIssues
	}
	return m.Issues
}

func (m *Model) ExitSearch() {
	m.SearchMode = false
	m.Query.SetValue("")
	m.applyFilter()
}

func (m *Model) applyFilter() {
	if m.Query.Value() == "" {
		m.FilteredIssues = m.Issues
		return
	}
	q := strings.ToLower(m.Query.Value())
	filtered := make([]repo.Issue, 0)
	for _, issue := range m.Issues {
		if strings.Contains(strings.ToLower(issue.Title), q) ||
			strings.Contains(strings.ToLower(issue.Body), q) ||
			fmt.Sprintf("%d", issue.Number) == q {
			filtered = append(filtered, issue)
		}
	}
	m.FilteredIssues = filtered
	m.Cursor = 0
	m.Offset = 0
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render(fmt.Sprintf(" Issues: %s ", m.Repo)))
	b.WriteString("\n\n")

	items := m.visibleItems()

	if m.SearchMode {
		b.WriteString(m.Query.View())
		b.WriteString("\n")
	}

	if len(items) == 0 {
		if m.SearchMode && m.Query.Value() != "" {
			b.WriteString("No matches.\n")
		} else {
			b.WriteString("No open issues found.\n")
		}
		b.WriteString(infoStyle.Render("q: quit"))
		return b.String()
	}

	visible := m.Height - 3
	if m.SearchMode {
		visible--
	}
	if visible < 1 {
		visible = 1
	}

	end := m.Offset + visible
	if end > len(items) {
		end = len(items)
	}

	for i := m.Offset; i < end; i++ {
		b.WriteString(renderLine(items[i], i == m.Cursor))
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	if m.SearchMode {
		b.WriteString(infoStyle.Render("↑/↓ or k/j: navigate  •  enter: view  •  esc: clear  •  q: quit"))
	} else {
		b.WriteString(infoStyle.Render("↑/↓ or k/j: navigate  •  /: search  •  o: open  •  enter: view  •  q: quit"))
	}
	return b.String()
}

func renderLine(issue repo.Issue, selected bool) string {
	prefix := "  "
	style := issueStyle
	if selected {
		prefix = "> "
		style = selectedStyle
	}

	stateTag := openColor.Render(issue.State)
	if issue.State != "OPEN" {
		stateTag = closedColor.Render(issue.State)
	}

	line := fmt.Sprintf("%s#%d [%s] %s", prefix, issue.Number, stateTag, issue.Title)
	return style.Render(line)
}

type SelectMsg struct{ Issue repo.Issue }
type OpenMsg struct{ Issue repo.Issue }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.SearchMode {
		return m.handleSearchMode(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case "enter":
			items := m.visibleItems()
			if m.Cursor < len(items) {
				return m, func() tea.Msg {
					return SelectMsg{Issue: items[m.Cursor]}
				}
			}
		case "o":
			items := m.visibleItems()
			if m.Cursor < len(items) {
				return m, func() tea.Msg {
					return OpenMsg{Issue: items[m.Cursor]}
				}
			}
		case "/":
			m.SearchMode = true
			m.Query.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m Model) handleSearchMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.SearchMode = false
			m.Query.SetValue("")
			m.applyFilter()
			return m, nil
		case "enter":
			items := m.visibleItems()
			if m.Cursor < len(items) {
				return m, func() tea.Msg {
					return SelectMsg{Issue: items[m.Cursor]}
				}
			}
			m.SearchMode = false
			return m, nil
		case "up":
			m.moveCursor(-1)
			return m, nil
		case "down":
			m.moveCursor(1)
			return m, nil
		case "o":
			items := m.visibleItems()
			if m.Cursor < len(items) {
				return m, func() tea.Msg {
					return OpenMsg{Issue: items[m.Cursor]}
				}
			}
		}
	}

	oldQuery := m.Query.Value()
	var cmd tea.Cmd
	m.Query, cmd = m.Query.Update(msg)
	if m.Query.Value() != oldQuery {
		m.applyFilter()
	}
	return m, cmd
}
