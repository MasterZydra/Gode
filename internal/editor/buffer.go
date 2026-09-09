package editor

import "strings"

type TextBuffer struct {
	text   []rune
	cursor int
}

func NewTextBuffer(text string) *TextBuffer {
	return &TextBuffer{text: []rune(text)}
}

func (b *TextBuffer) String() string {
	return string(b.text)
}

func (b *TextBuffer) Cursor() int {
	return b.cursor
}

func (b *TextBuffer) SetCursor(cursor int) {
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(b.text) {
		cursor = len(b.text)
	}
	b.cursor = cursor
}

func (b *TextBuffer) SetText(text string) {
	b.text = []rune(text)
	b.cursor = 0
}

func (b *TextBuffer) Insert(r rune) {
	b.text = append(b.text, 0)
	copy(b.text[b.cursor+1:], b.text[b.cursor:])
	b.text[b.cursor] = r
	b.cursor++
}

func (b *TextBuffer) Backspace() {
	if b.cursor == 0 {
		return
	}
	b.text = append(b.text[:b.cursor-1], b.text[b.cursor:]...)
	b.cursor--
}

func (b *TextBuffer) Delete() {
	if b.cursor >= len(b.text) {
		return
	}
	b.text = append(b.text[:b.cursor], b.text[b.cursor+1:]...)
}

func (b *TextBuffer) ReplaceRange(start, end int, replacement []rune) {
	if start < 0 {
		start = 0
	}
	if end > len(b.text) {
		end = len(b.text)
	}
	if start > end {
		start, end = end, start
	}
	updated := make([]rune, 0, len(b.text)-(end-start)+len(replacement))
	updated = append(updated, b.text[:start]...)
	updated = append(updated, replacement...)
	updated = append(updated, b.text[end:]...)
	b.text = updated
	b.SetCursor(start + len(replacement))
}

func (b *TextBuffer) MoveLeft() {
	if b.cursor > 0 {
		b.cursor--
	}
}

func (b *TextBuffer) MoveRight() {
	if b.cursor < len(b.text) {
		b.cursor++
	}
}

func (b *TextBuffer) MoveHome() {
	lineStart, _ := b.LineBounds()
	b.cursor = lineStart
}

func (b *TextBuffer) MoveEnd() {
	_, lineEnd := b.LineBounds()
	b.cursor = lineEnd
}

func (b *TextBuffer) MoveUp() {
	line, column := b.CursorLineColumn()
	if line == 0 {
		return
	}
	b.cursor = b.offsetForLineColumn(line-1, column)
}

func (b *TextBuffer) MoveDown() {
	line, column := b.CursorLineColumn()
	lines := strings.Split(b.String(), "\n")
	if line >= len(lines)-1 {
		return
	}
	b.cursor = b.offsetForLineColumn(line+1, column)
}

func (b *TextBuffer) CursorLineColumn() (int, int) {
	line := 0
	column := 0
	for i, r := range b.text {
		if i == b.cursor {
			return line, column
		}
		if r == '\n' {
			line++
			column = 0
		} else {
			column++
		}
	}
	return line, column
}

func (b *TextBuffer) LineBounds() (int, int) {
	start := b.cursor
	for start > 0 && b.text[start-1] != '\n' {
		start--
	}
	end := b.cursor
	for end < len(b.text) && b.text[end] != '\n' {
		end++
	}
	return start, end
}

func (b *TextBuffer) Lines() []string {
	return strings.Split(b.String(), "\n")
}

func (b *TextBuffer) offsetForLineColumn(line, column int) int {
	currentLine := 0
	offset := 0
	for offset < len(b.text) && currentLine < line {
		if b.text[offset] == '\n' {
			currentLine++
		}
		offset++
	}
	lineEnd := offset
	for lineEnd < len(b.text) && b.text[lineEnd] != '\n' {
		lineEnd++
	}
	if offset+column > lineEnd {
		return lineEnd
	}
	return offset + column
}
