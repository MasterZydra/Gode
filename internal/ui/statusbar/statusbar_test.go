package statusbar

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestStatusBarShowsBranchAndAheadBehind(t *testing.T) {
	test.NewTempApp(t)
	root := t.TempDir()
	runStatusBarGit(t, root, "init", "-b", "main")
	runStatusBarGit(t, root, "config", "user.email", "test@example.com")
	runStatusBarGit(t, root, "config", "user.name", "Test User")
	writeStatusBarFile(t, filepath.Join(root, "file.txt"))
	runStatusBarGit(t, root, "add", "file.txt")
	runStatusBarGit(t, root, "commit", "-m", "initial")

	bar := NewStatusBar(root)
	if bar.label.Text != "main · 0 ahead · 0 behind" {
		t.Fatalf("got status text %q", bar.label.Text)
	}
}

func runStatusBarGit(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}

func writeStatusBarFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
}
