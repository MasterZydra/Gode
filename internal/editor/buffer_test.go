package editor

import "testing"

func TestTextBufferEditing(t *testing.T) {
	buffer := NewTextBuffer("ab\ncd")
	buffer.SetCursor(2)
	buffer.Insert('X')
	buffer.Backspace()
	buffer.MoveDown()
	buffer.MoveEnd()
	buffer.Delete()

	if got := buffer.String(); got != "ab\ncd" {
		t.Fatalf("unexpected buffer text %q", got)
	}
}

func TestTextBufferCursorLineColumn(t *testing.T) {
	buffer := NewTextBuffer("one\ntwo")
	buffer.SetCursor(5)

	line, column := buffer.CursorLineColumn()
	if line != 1 || column != 1 {
		t.Fatalf("cursor = (%d, %d), want (1, 1)", line, column)
	}
}
