package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseStatus(t *testing.T) {
	status, err := ParseStatus(" M changed file.go\nM  staged.go\n?? new file.txt\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Staged) != 1 || status.Staged[0].Path != "staged.go" {
		t.Fatalf("got staged %#v", status.Staged)
	}
	if len(status.Unstaged) != 2 || status.Unstaged[0].Path != "changed file.go" || status.Unstaged[1].Path != "new file.txt" {
		t.Fatalf("got unstaged %#v", status.Unstaged)
	}
}

func TestRepositoryStageUnstageAndCommit(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "Test User")

	filePath := filepath.Join(root, "file.txt")
	writeFile(t, filePath, "initial")
	runGit(t, root, "add", "file.txt")
	runGit(t, root, "commit", "-m", "initial")

	writeFile(t, filePath, "changed")
	repository := NewRepository(root)
	status, err := repository.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Unstaged) != 1 || status.Unstaged[0].Path != "file.txt" {
		t.Fatalf("got status %#v", status)
	}

	if err := repository.Stage("file.txt"); err != nil {
		t.Fatal(err)
	}
	status, err = repository.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Staged) != 1 || len(status.Unstaged) != 0 {
		t.Fatalf("got staged status %#v", status)
	}

	if err := repository.Unstage("file.txt"); err != nil {
		t.Fatal(err)
	}
	status, err = repository.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Staged) != 0 || len(status.Unstaged) != 1 {
		t.Fatalf("got unstaged status %#v", status)
	}

	if err := repository.Stage("file.txt"); err != nil {
		t.Fatal(err)
	}
	if err := repository.Commit("update file"); err != nil {
		t.Fatal(err)
	}
	status, err = repository.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Staged) != 0 || len(status.Unstaged) != 0 {
		t.Fatalf("got status after commit %#v", status)
	}
}

func TestRepositoryStatusReportsUntrackedFile(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	writeFile(t, filepath.Join(root, "new file.txt"), "new")

	status, err := NewRepository(root).Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Unstaged) != 1 || status.Unstaged[0].Path != "new file.txt" {
		t.Fatalf("got status %#v", status)
	}
}

func TestRepositoryInfoReportsBranchAndAheadBehind(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init", "-b", "main")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "Test User")
	writeFile(t, filepath.Join(root, "file.txt"), "initial")
	runGit(t, root, "add", "file.txt")
	runGit(t, root, "commit", "-m", "initial")

	info, err := NewRepository(root).Info()
	if err != nil {
		t.Fatal(err)
	}
	if info.Branch != "main" || info.Ahead != 0 || info.Behind != 0 {
		t.Fatalf("got info %#v", info)
	}
}

func TestRepositoryPushAndPull(t *testing.T) {
	remote := t.TempDir()
	runGit(t, remote, "init", "--bare", "-b", "main")

	root := t.TempDir()
	runGit(t, root, "init", "-b", "main")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "Test User")
	runGit(t, root, "remote", "add", "origin", remote)
	filePath := filepath.Join(root, "file.txt")
	writeFile(t, filePath, "initial")
	runGit(t, root, "add", "file.txt")
	runGit(t, root, "commit", "-m", "initial")
	runGit(t, root, "push", "--set-upstream", "origin", "main")

	repository := NewRepository(root)

	writeFile(t, filePath, "pushed")
	runGit(t, root, "add", "file.txt")
	runGit(t, root, "commit", "-m", "pushed")
	if err := repository.Push(); err != nil {
		t.Fatal(err)
	}

	clone := t.TempDir()
	runGit(t, clone, "clone", "--branch", "main", remote, ".")
	runGit(t, clone, "config", "user.email", "test@example.com")
	runGit(t, clone, "config", "user.name", "Test User")
	cloneFile := filepath.Join(clone, "file.txt")
	writeFile(t, cloneFile, "pulled")
	runGit(t, clone, "add", "file.txt")
	runGit(t, clone, "commit", "-m", "pulled")
	runGit(t, clone, "push")

	if err := repository.Pull(); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "pulled" {
		t.Fatalf("got file contents %q, want pulled", contents)
	}
}

func TestRepositoryDiff(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "Test User")
	path := filepath.Join(root, "file.txt")
	writeFile(t, path, "before\n")
	runGit(t, root, "add", "file.txt")
	runGit(t, root, "commit", "-m", "initial")
	writeFile(t, path, "after\n")

	diff, err := NewRepository(root).Diff("file.txt", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(diff, "-before") || !strings.Contains(diff, "+after") {
		t.Fatalf("diff does not contain expected changes: %q", diff)
	}
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}

func writeFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
