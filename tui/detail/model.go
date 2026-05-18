package detail

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ghi/render"
	"ghi/repo"
)

var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)
)

type Model struct {
	Issue  repo.Issue
	VP     viewport.Model
	Width  int
	Height int
}

func New(w, h int) Model {
	vp := viewport.New(w, h-3)
	vp.YPosition = 1
	return Model{VP: vp, Width: w, Height: h}
}

func (m *Model) UpdateSize(w, h int) {
	m.Width = w
	m.Height = h
	m.VP.Width = w
	if h > 3 {
		m.VP.Height = h - 3
	} else {
		m.VP.Height = 1
	}
}

func (m *Model) SetIssue(issue repo.Issue, renderer render.Renderer) {
	m.Issue = issue
	m.VP.SetContent(renderer.Render(issue, m.Width-4))
	m.VP.GotoTop()
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render(fmt.Sprintf(" #%d ", m.Issue.Number)))
	b.WriteString("\n")
	b.WriteString(m.VP.View())
	b.WriteString("\n")
	b.WriteString(infoStyle.Render("↑/↓ or k/j: scroll  •  g/G: top/bottom  •  o: open  •  esc/b: back  •  q: quit"))
	return b.String()
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.VP.LineUp(1)
		case "down", "j":
			m.VP.LineDown(1)
		case "g":
			m.VP.GotoTop()
		case "G":
			m.VP.GotoBottom()
		}
	}
	newVp, cmd := m.VP.Update(msg)
	m.VP = newVp
	return m, cmd
}
