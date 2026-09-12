package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestGitViewCommitButtonRequiresMessage(t *testing.T) {
	test.NewTempApp(t)
	view := NewGitView(t.TempDir(), nil, nil)

	if !view.commit.Disabled() {
		t.Fatal("commit button is enabled for an empty message")
	}

	view.message.SetText("commit message")
	if view.commit.Disabled() {
		t.Fatal("commit button is disabled for a non-empty message")
	}

	view.message.SetText("   ")
	if !view.commit.Disabled() {
		t.Fatal("commit button is enabled for a whitespace-only message")
	}
}

func TestGitViewHidesEmptyStagedSection(t *testing.T) {
	test.NewTempApp(t)
	root := t.TempDir()
	runGitCommand(t, root, "init")
	writeGitTestFile(t, filepath.Join(root, "file.txt"), "content")

	view := NewGitView(root, nil, nil)
	if len(view.stagedSection.Objects) != 0 {
		t.Fatalf("staged section has %d objects, want 0", len(view.stagedSection.Objects))
	}
	if len(view.changesSection.Objects) != 2 {
		t.Fatalf("changes section has %d objects, want heading and file", len(view.changesSection.Objects))
	}

	runGitCommand(t, root, "add", "file.txt")
	view.Refresh()
	if len(view.stagedSection.Objects) != 2 {
		t.Fatalf("staged section has %d objects, want heading and file", len(view.stagedSection.Objects))
	}
}

func TestGitViewRefreshButtonReloadsStatus(t *testing.T) {
	test.NewTempApp(t)
	root := t.TempDir()
	runGitCommand(t, root, "init")
	view := NewGitView(root, nil, nil)

	writeGitTestFile(t, filepath.Join(root, "new.txt"), "content")
	if len(view.changesSection.Objects) != 2 {
		t.Fatalf("changes section has %d objects before refresh, want heading and no changes", len(view.changesSection.Objects))
	}

	test.Tap(view.refresh)
	if len(view.changesSection.Objects) != 2 {
		t.Fatalf("changes section has %d objects after refresh, want heading and file", len(view.changesSection.Objects))
	}
}

func runGitCommand(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}

func writeGitTestFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
