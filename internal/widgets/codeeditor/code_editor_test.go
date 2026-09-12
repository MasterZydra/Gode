package codeeditor

import (
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func TestMain(main *testing.M) {
	app.NewWithID("com.gode.editor.tests")
	os.Exit(main.Run())
}

func TestDoubleClickSelectsWord(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("hello world")
	position := fyne.NewPos(editor.lineNumberWidth()+editor.charWidth()*7, 0)

	editor.DoubleTapped(&fyne.PointEvent{Position: position})

	if got := editor.selectedText(); got != "world" {
		t.Fatalf("selected text = %q, want %q", got, "world")
	}
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

func TestCodeEditorReadOnlyRejectsInput(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("before")
	editor.SetReadOnly(true)
	editor.focused = true
	editor.TypedRune('x')

	if got := editor.Text(); got != "before" {
		t.Fatalf("read-only editor changed text to %q", got)
	}
}

func TestCodeEditorInsertsTab(t *testing.T) {
	editor := NewCodeEditor()
	editor.focused = true
	if !editor.AcceptsTab() {
		t.Fatal("editor does not accept Tab")
	}

	editor.TypedKey(&fyne.KeyEvent{Name: fyne.KeyTab})

	if got := editor.Text(); got != "\t" {
		t.Fatalf("text = %q, want a tab", got)
	}
}
