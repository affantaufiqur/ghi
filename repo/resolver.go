package repo

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var ErrNotGitHub = fmt.Errorf("Not a git repo or not connecting to a github remote")

type Resolver interface {
	Slug() (string, error)
}

type GitRemote struct{}

func (GitRemote) Slug() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) != "true" {
		return "", ErrNotGitHub
	}

	out, err = exec.Command("git", "remote", "-v").CombinedOutput()
	if err != nil {
		return "", ErrNotGitHub
	}

	re := regexp.MustCompile(`github\.com[:/]([^/]+)/([^/\s]+)`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 3 {
		return "", ErrNotGitHub
	}

	owner := matches[1]
	repo := strings.TrimSuffix(matches[2], ".git")
	return fmt.Sprintf("%s/%s", owner, repo), nil
}
