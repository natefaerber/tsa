package cmd

import (
	"fmt"
	"os/exec"
)

// GitClient abstracts git operations for testability.
type GitClient interface {
	StagedDiff(files []string) (string, error)
}

// execGitClient shells out to the system git binary.
type execGitClient struct{}

func (execGitClient) StagedDiff(files []string) (string, error) {
	args := []string{"diff", "--cached", "--diff-filter=ACM"}
	if len(files) > 0 {
		args = append(args, "--")
		args = append(args, files...)
	}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", fmt.Errorf("running git diff: %w", err)
	}
	return string(out), nil
}

// git is the active client, replaced in tests.
var git GitClient = execGitClient{}
