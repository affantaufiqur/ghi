package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"ghi/repo"
)

type Renderer interface {
	Render(issue repo.Issue, width int) string
}

type Glamour struct {
	r      *glamour.TermRenderer
	width  int
}

func NewGlamour() *Glamour {
	return &Glamour{}
}

func (g *Glamour) SetWidth(width int) {
	if g.width == width && g.r != nil {
		return
	}
	g.width = width
	g.r, _ = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
}

func (g *Glamour) Render(issue repo.Issue, width int) string {
	g.SetWidth(width)

	labelStr := joinNames(issue.Labels, ", ", "none")
	assigneeStr := joinLogins(issue.Assignees, ", ", "unassigned")

	md := fmt.Sprintf(
		"# #%d %s\n\n**State:** %s  |  **Author:** @%s  |  **Created:** %s\n\n**Labels:** %s  |  **Assignees:** %s\n\n---\n\n%s",
		issue.Number,
		issue.Title,
		issue.State,
		issue.Author.Login,
		issue.CreatedAt,
		labelStr,
		assigneeStr,
		issue.Body,
	)

	if g.r == nil {
		return issue.Body
	}
	rendered, err := g.r.Render(md)
	if err != nil {
		return issue.Body
	}
	return rendered
}

type Plain struct{}

func (Plain) Render(issue repo.Issue, width int) string {
	open := lipgloss.NewStyle().Foreground(lipgloss.Color("#3FB950")).Render("OPEN")
	closed := lipgloss.NewStyle().Foreground(lipgloss.Color("#F85149")).Render("CLOSED")
	state := open
	if issue.State != "OPEN" {
		state = closed
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("#%d [%s] %s\n", issue.Number, state, issue.Title))
	b.WriteString(fmt.Sprintf("Author: @%s  Created: %s\n", issue.Author.Login, issue.CreatedAt))
	b.WriteString(fmt.Sprintf("Labels: %s  Assignees: %s\n", joinNames(issue.Labels, ", ", "none"), joinLogins(issue.Assignees, ", ", "unassigned")))
	b.WriteString("---\n")
	b.WriteString(issue.Body)
	return b.String()
}

func joinNames(labels []repo.Label, sep, fallback string) string {
	if len(labels) == 0 {
		return fallback
	}
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return strings.Join(names, sep)
}

func joinLogins(authors []repo.Author, sep, fallback string) string {
	if len(authors) == 0 {
		return fallback
	}
	logins := make([]string, len(authors))
	for i, a := range authors {
		logins[i] = "@" + a.Login
	}
	return strings.Join(logins, sep)
}
