package editor

import "testing"

func TestPasteCopiedLineBelowCurrentLine(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("one\ntwo")
	editor.cursors = []cursorState{{position: 1, anchor: 1}}
	editor.lineClipboardText = "one"
	editor.lineClipboard = true

	editor.pasteLinesBelow()

	if got := editor.Text(); got != "one\none\ntwo" {
		t.Fatalf("got %q, want %q", got, "one\none\ntwo")
	}
}

func TestPasteCopiedFinalLineBelowCurrentLine(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("one\ntwo")
	editor.cursors = []cursorState{{position: len([]rune(editor.Text())), anchor: len([]rune(editor.Text()))}}
	editor.lineClipboardText = "two"
	editor.lineClipboard = true

	editor.pasteLinesBelow()

	if got := editor.Text(); got != "one\ntwo\ntwo" {
		t.Fatalf("got %q, want %q", got, "one\ntwo\ntwo")
	}
}
