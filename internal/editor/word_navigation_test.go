package editor

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func TestWordNavigationUsesPunctuationAndWhitespace(t *testing.T) {
	editor := NewCodeEditor()
	editor.Buffer.SetText("class.method")
	editor.cursors = []cursorState{{position: len([]rune(editor.Text())), anchor: len([]rune(editor.Text()))}}

	editor.moveCursorsWord(true, false)
	if got := editor.cursors[0].position; got != 6 {
		t.Fatalf("first word-left position = %d, want 6", got)
	}
	editor.cursors[0].anchor = editor.cursors[0].position
	editor.moveCursorsWord(true, false)
	if got := editor.cursors[0].position; got != 5 {
		t.Fatalf("second word-left position = %d, want 5", got)
	}
	editor.cursors[0].anchor = editor.cursors[0].position
	editor.moveCursorsWord(true, false)
	if got := editor.cursors[0].position; got != 0 {
		t.Fatalf("third word-left position = %d, want 0", got)
	}
	editor.cursors[0].anchor = editor.cursors[0].position
	editor.moveCursorsWord(false, false)
	if got := editor.cursors[0].position; got != 5 {
		t.Fatalf("word-right position = %d, want 5", got)
	}
	editor.cursors[0].anchor = editor.cursors[0].position
	editor.moveCursorsWord(false, false)
	if got := editor.cursors[0].position; got != 6 {
		t.Fatalf("second word-right position = %d, want 6", got)
	}
}

func TestCustomShortcutWordNavigation(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("one two")
	editor.cursors = []cursorState{{position: len([]rune(editor.Text())), anchor: len([]rune(editor.Text()))}}

	editor.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyLeft, Modifier: fyne.KeyModifierControl})
	if got := editor.cursors[0].position; got != 4 {
		t.Fatalf("Ctrl+Left position = %d, want 4", got)
	}
}

func TestRepeatedShiftWordNavigationExtendsSelection(t *testing.T) {
	editor := NewCodeEditor()
	editor.SetText("class.method")
	editor.cursors = []cursorState{{position: 0, anchor: 0}}

	editor.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift})
	editor.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift})

	if got := editor.selectedText(); got != "class." {
		t.Fatalf("selected text = %q, want %q", got, "class.")
	}
}
