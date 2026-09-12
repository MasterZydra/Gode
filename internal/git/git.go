package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Change struct {
	Path string
}

type Status struct {
	Staged   []Change
	Unstaged []Change
}

type Info struct {
	Branch string
	Ahead  int
	Behind int
}

type Repository struct {
	root string
}

func NewRepository(root string) *Repository {
	return &Repository{root: root}
}

func (r *Repository) Status() (Status, error) {
	output, err := r.run("status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return Status{}, err
	}
	return ParseStatus(string(output))
}

func (r *Repository) Info() (Info, error) {
	branchOutput, err := r.run("branch", "--show-current")
	if err != nil {
		return Info{}, err
	}

	info := Info{Branch: strings.TrimSpace(string(branchOutput))}
	if info.Branch == "" {
		info.Branch = "HEAD"
	}

	counts, err := r.run("rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if err != nil {
		return info, nil
	}
	if _, err := fmt.Sscanf(string(counts), "%d %d", &info.Behind, &info.Ahead); err != nil {
		return Info{}, fmt.Errorf("invalid git ahead/behind output %q: %w", counts, err)
	}
	return info, nil
}

func (r *Repository) Stage(path string) error {
	_, err := r.run("add", "--", path)
	return err
}

func (r *Repository) Unstage(path string) error {
	_, err := r.run("restore", "--staged", "--", path)
	return err
}

func (r *Repository) Commit(message string) error {
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("commit message must not be empty")
	}
	_, err := r.run("commit", "-m", message)
	return err
}

func ParseStatus(output string) (Status, error) {
	var status Status
	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		if line == "" {
			continue
		}
		if len(line) < 3 {
			return Status{}, fmt.Errorf("invalid git status line %q", line)
		}

		indexStatus := line[0]
		worktreeStatus := line[1]
		path := line[3:]
		if path == "" {
			return Status{}, fmt.Errorf("missing path in git status line %q", line)
		}
		if strings.HasPrefix(path, "\"") {
			unquotedPath, err := strconv.Unquote(path)
			if err != nil {
				return Status{}, fmt.Errorf("invalid quoted path in git status line %q: %w", line, err)
			}
			path = unquotedPath
		}

		if indexStatus != ' ' && indexStatus != '?' {
			status.Staged = append(status.Staged, Change{Path: path})
		}
		if worktreeStatus != ' ' || (indexStatus == '?' && worktreeStatus == '?') {
			status.Unstaged = append(status.Unstaged, Change{Path: path})
		}
	}
	return status, nil
}

func (r *Repository) run(args ...string) ([]byte, error) {
	if strings.TrimSpace(r.root) == "" {
		return nil, fmt.Errorf("repository path is empty")
	}
	commandArgs := append([]string{"-C", r.root}, args...)
	command := exec.Command("git", commandArgs...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), message)
	}
	return output, nil
}
