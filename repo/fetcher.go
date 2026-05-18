package repo

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Fetcher interface {
	Fetch() ([]Issue, error)
}

type Issue struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	State     string   `json:"state"`
	Author    Author   `json:"author"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Body      string   `json:"body"`
	Labels    []Label  `json:"labels"`
	Assignees []Author `json:"assignees"`
}

type Author struct {
	Login string `json:"login"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type GHCLI struct{}

func (GHCLI) Fetch() ([]Issue, error) {
	cmd := exec.Command(
		"gh", "issue", "list",
		"--json", "number,title,author,state,createdAt,updatedAt,body,labels,assignees",
		"--limit", "100",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %s", string(out))
	}

	var issues []Issue
	if err := json.Unmarshal(out, &issues); err != nil {
		return nil, fmt.Errorf("failed to parse issues: %w", err)
	}
	return issues, nil
}
