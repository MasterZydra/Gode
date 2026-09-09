package editor

import (
	"os"
	"testing"

	"fyne.io/fyne/v2/app"
)

func TestMain(main *testing.M) {
	app.NewWithID("com.gode.editor.tests")
	os.Exit(main.Run())
}

func TestCodeEditorReplacesSelections(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("one two")
	editor.cursors = []cursorState{
		{position: 3, anchor: 0},
		{position: 7, anchor: 4},
	}

	editor.replaceSelections([]rune("X"))

	if got := editor.Text(); got != "X X" {
		t.Fatalf("got %q, want %q", got, "X X")
	}
}

func TestCodeEditorMultipleCursorsInsert(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("aa\nbb")
	editor.cursors = []cursorState{
		{position: 0, anchor: 0},
		{position: 3, anchor: 3},
	}

	editor.replaceSelections([]rune("X"))

	if got := editor.Text(); got != "Xaa\nXbb" {
		t.Fatalf("got %q, want %q", got, "Xaa\nXbb")
	}
}
