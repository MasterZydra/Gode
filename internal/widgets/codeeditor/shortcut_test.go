package codeeditor

import "testing"

func TestSelectAllSelectsWholeFile(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("one\ntwo")

	editor.selectAll()

	if got := editor.selectedText(); got != "one\ntwo" {
		t.Fatalf("selected text = %q, want %q", got, "one\ntwo")
	}
}

func TestCutCurrentLineRemovesLine(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("one\ntwo\nthree")
	editor.cursors = []cursorState{{position: 5, anchor: 5}}
	editor.lineClipboardText = editor.currentLinesText()
	editor.lineClipboard = true

	editor.cutCurrentLines()

	if got := editor.Text(); got != "one\nthree" {
		t.Fatalf("text = %q, want %q", got, "one\nthree")
	}
	if editor.lineClipboardText != "two" {
		t.Fatalf("line clipboard = %q, want %q", editor.lineClipboardText, "two")
	}
}
