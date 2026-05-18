# ghi

Terminal UI for browsing GitHub issues in whatever repo you're in right now. Uses Bubble Tea.

> Built entirely by vibe coding.

Run `ghi` inside a git repo with a GitHub remote. It reads your remote, runs `gh issue list`, and shows the results.

## What you get

- No setup — reads the repo from `git remote -v`
- Search with `/` — filters by title, body, or issue number
- Markdown rendered with glamour
- `o` to open in browser
- `j`/`k` navigation

## Needs

- `git`
- [`gh`](https://cli.github.com/) (GitHub CLI)

## Install

```bash
make build
make install   # copies to ~/.local/bin
```

Or just:

```bash
go run ./cmd/ghi
```

## Usage

```bash
cd my-project
ghi
```

### Keys

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Move |
| `Enter` | View issue |
| `/` | Search |
| `o` | Open in browser |
| `Esc` or `b` | Back to list |
| `g` / `G` | Top / bottom |
| `q` or `Ctrl+c` | Quit |

## Built with

- Go 1.24
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Glamour](https://github.com/charmbracelet/glamour)

## License

MIT
