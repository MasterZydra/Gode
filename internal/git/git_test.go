package git

import (
	"os"
	"os/exec"
	"path/filepath"
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
