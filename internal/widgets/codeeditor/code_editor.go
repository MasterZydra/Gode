package codeeditor

import (
	"fmt"
	"gode/internal/highlighter"
	"image/color"
	"sort"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	keywordColor = color.NRGBA{R: 0, G: 102, B: 204, A: 255}
	commentColor = color.NRGBA{R: 80, G: 180, B: 90, A: 255}
	stringColor  = color.NRGBA{R: 190, G: 100, B: 20, A: 255}
	numberColor  = color.NRGBA{R: 80, G: 180, B: 90, A: 255}
	cursorColor  = color.NRGBA{R: 40, G: 120, B: 220, A: 255}
)

type CodeEditor struct {
	widget.BaseWidget
	Buffer            *TextBuffer
	FileName          string
	OnChanged         func(string)
	OnCursorChanged   func()
	focused           bool
	highlighter       highlighter.Highlighter
	lines             []highlighter.HighlightedLine
	cursors           []cursorState
	dragging          bool
	shiftDown         bool
	lineClipboard     bool
	lineClipboardText string
	readOnly          bool
}

type cursorState struct {
	position int
	anchor   int
}

var (
	_ fyne.Focusable      = (*CodeEditor)(nil)
	_ fyne.Tabbable       = (*CodeEditor)(nil)
	_ desktop.Keyable     = (*CodeEditor)(nil)
	_ desktop.Mouseable   = (*CodeEditor)(nil)
	_ fyne.Draggable      = (*CodeEditor)(nil)
	_ fyne.DoubleTappable = (*CodeEditor)(nil)
)

func NewCodeEditor() *CodeEditor {
	e := &CodeEditor{Buffer: NewTextBuffer(""), cursors: []cursorState{{}}}
	e.ExtendBaseWidget(e)
	e.refreshHighlighting()
	return e
}

func (e *CodeEditor) SetText(text string) {
	e.Buffer.SetText(text)
	e.cursors = []cursorState{{}}
	e.refreshHighlighting()
	e.Refresh()
}

func (e *CodeEditor) Text() string {
	return e.Buffer.String()
}

func (e *CodeEditor) AcceptsTab() bool {
	return true
}

func (e *CodeEditor) SetReadOnly(readOnly bool) {
	e.readOnly = readOnly
}

func (e *CodeEditor) CursorBounds() (fyne.Position, fyne.Size) {
	if len(e.cursors) == 0 {
		return fyne.Position{}, fyne.Size{}
	}

	position := e.cursors[0].position
	e.Buffer.SetCursor(position)
	line, rawColumn := e.Buffer.CursorLineColumn()
	lineText := e.Buffer.Lines()[line]
	column := highlighter.VisualColumn(string([]rune(lineText)[:rawColumn]))
	return fyne.NewPos(e.lineNumberWidth()+float32(column)*e.charWidth(), float32(line)*e.lineHeight()), fyne.NewSize(2, e.lineHeight())
}

func (e *CodeEditor) SetHighlighter(syntaxHighlighter highlighter.Highlighter) {
	e.highlighter = syntaxHighlighter
	e.refreshHighlighting()
	e.Refresh()
}

func (e *CodeEditor) refreshHighlighting() {
	if e.highlighter != nil {
		e.lines = e.highlighter.Highlight(e.Buffer.String())
	} else {
		e.lines = highlighter.PlainLines(e.Buffer.String())
	}
}

func (e *CodeEditor) changed() {
	e.syncPrimaryCursor()
	e.refreshHighlighting()
	e.Refresh()
	if e.OnChanged != nil {
		e.OnChanged(e.Text())
	}
}

func (e *CodeEditor) TypedRune(r rune) {
	if !e.focused || e.readOnly {
		return
	}
	e.lineClipboard = false
	e.replaceSelections([]rune{r})
	e.changed()
}

func (e *CodeEditor) TypedKey(key *fyne.KeyEvent) {
	if !e.focused || e.readOnly {
		return
	}
	oldCursors := append([]cursorState(nil), e.cursors...)
	controlDown := e.controlDown()
	switch key.Name {
	case fyne.KeyBackspace:
		e.lineClipboard = false
		e.deleteFromCursors(false)
	case fyne.KeyDelete:
		e.lineClipboard = false
		e.deleteFromCursors(true)
	case fyne.KeyLeft:
		if controlDown {
			e.moveCursorsWord(true, e.shiftDown)
		} else {
			e.moveCursors(func(position int) int {
				if position > 0 {
					return position - 1
				}
				return position
			})
		}
	case fyne.KeyRight:
		if controlDown {
			e.moveCursorsWord(false, e.shiftDown)
		} else {
			e.moveCursors(func(position int) int {
				if position < len([]rune(e.Text())) {
					return position + 1
				}
				return position
			})
		}
	case fyne.KeyUp:
		e.moveCursorsLine(true)
	case fyne.KeyDown:
		e.moveCursorsLine(false)
	case fyne.KeyHome:
		e.moveCursors(func(position int) int { e.Buffer.SetCursor(position); start, _ := e.Buffer.LineBounds(); return start })
	case fyne.KeyEnd:
		e.moveCursors(func(position int) int { e.Buffer.SetCursor(position); _, end := e.Buffer.LineBounds(); return end })
	case fyne.KeyReturn, fyne.KeyEnter:
		e.lineClipboard = false
		e.replaceSelections([]rune{'\n'})
	case fyne.KeyTab:
		e.lineClipboard = false
		e.replaceSelections([]rune{'\t'})
	default:
		return
	}
	if e.shiftDown {
		for index := range e.cursors {
			e.cursors[index].anchor = oldCursors[index].anchor
		}
	} else {
		for index := range e.cursors {
			e.cursors[index].anchor = e.cursors[index].position
		}
	}
	e.changed()
}

func (e *CodeEditor) controlDown() bool {
	desktopDriver, ok := fyne.CurrentApp().Driver().(desktop.Driver)
	return ok && desktopDriver.CurrentKeyModifiers()&fyne.KeyModifierControl != 0
}

func (e *CodeEditor) moveCursorsWord(left, extendSelection bool) {
	runes := []rune(e.Text())
	for index := range e.cursors {
		position := e.cursors[index].position
		if !extendSelection && e.cursors[index].anchor != position {
			if left {
				if e.cursors[index].anchor < position {
					position = e.cursors[index].anchor
				}
			} else if e.cursors[index].anchor > position {
				position = e.cursors[index].anchor
			}
			e.cursors[index].position = position
			continue
		}
		if left {
			if position > 0 && isWordRune(runes[position-1]) {
				for position > 0 && isWordRune(runes[position-1]) {
					position--
				}
			} else if position > 0 {
				position--
			}
		} else {
			if position < len(runes) && isWordRune(runes[position]) {
				for position < len(runes) && isWordRune(runes[position]) {
					position++
				}
			} else if position < len(runes) {
				position++
			}
		}
		e.cursors[index].position = position
	}
}

func isWordRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r)
}

func (e *CodeEditor) KeyDown(key *fyne.KeyEvent) {
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		e.shiftDown = true
	}
}

func (e *CodeEditor) KeyUp(key *fyne.KeyEvent) {
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		e.shiftDown = false
	}
}

func (e *CodeEditor) TypedShortcut(shortcut fyne.Shortcut) {
	if e.readOnly {
		return
	}
	if custom, ok := shortcut.(*desktop.CustomShortcut); ok {
		if custom.Modifier&fyne.KeyModifierControl != 0 && (custom.KeyName == fyne.KeyLeft || custom.KeyName == fyne.KeyRight) {
			oldCursors := append([]cursorState(nil), e.cursors...)
			e.moveCursorsWord(custom.KeyName == fyne.KeyLeft, custom.Modifier&fyne.KeyModifierShift != 0)
			if custom.Modifier&fyne.KeyModifierShift != 0 {
				for index := range e.cursors {
					e.cursors[index].anchor = oldCursors[index].anchor
				}
			} else {
				for index := range e.cursors {
					e.cursors[index].anchor = e.cursors[index].position
				}
			}
			e.changed()
		}
		return
	}

	clipboard := fyne.CurrentApp().Clipboard()
	switch shortcut.(type) {
	case *fyne.ShortcutCopy:
		if selected := e.selectedText(); selected != "" {
			e.lineClipboard = false
			clipboard.SetContent(selected)
		} else {
			e.lineClipboardText = e.currentLinesText()
			e.lineClipboard = true
			clipboard.SetContent(e.lineClipboardText)
		}
		return
	case *fyne.ShortcutCut:
		if selected := e.selectedText(); selected != "" {
			e.lineClipboard = false
			clipboard.SetContent(selected)
			e.deleteSelections()
		} else {
			e.lineClipboardText = e.currentLinesText()
			e.lineClipboard = true
			clipboard.SetContent(e.lineClipboardText)
			e.cutCurrentLines()
		}
	case *fyne.ShortcutSelectAll:
		e.selectAll()
		return
	case *fyne.ShortcutPaste:
		if e.lineClipboard {
			e.pasteLinesBelow()
			e.lineClipboard = false
		} else {
			e.lineClipboard = false
			e.replaceSelections([]rune(clipboard.Content()))
		}
	default:
		return
	}
	e.changed()
}

func (e *CodeEditor) FocusGained() {
	e.focused = true
	e.Refresh()
}

func (e *CodeEditor) FocusLost() {
	e.focused = false
	e.Refresh()
}

func (e *CodeEditor) MouseDown(event *desktop.MouseEvent) {
	e.focused = true
	e.dragging = true
	e.shiftDown = event.Modifier&fyne.KeyModifierShift != 0
	if canvas := fyne.CurrentApp().Driver().CanvasForObject(e); canvas != nil {
		canvas.Focus(e)
	}
	position := e.cursorOffsetFromPosition(event.Position)
	if event.Modifier&fyne.KeyModifierAlt != 0 {
		e.cursors = append(e.cursors, cursorState{position: position, anchor: position})
	} else {
		if e.shiftDown {
			e.cursors[0].position = position
		} else {
			e.cursors = []cursorState{{position: position, anchor: position}}
		}
	}
	e.syncPrimaryCursor()
	e.Refresh()
}

func (e *CodeEditor) MouseUp(*desktop.MouseEvent)    { e.dragging = false }
func (e *CodeEditor) MouseMoved(*desktop.MouseEvent) {}
func (e *CodeEditor) MouseOut()                      {}
func (e *CodeEditor) MouseIn(*desktop.MouseEvent)    {}
func (e *CodeEditor) Dragged(event *fyne.DragEvent) {
	if !e.dragging {
		return
	}
	e.cursors[0].position = e.cursorOffsetFromPosition(event.Position)
	e.syncPrimaryCursor()
	e.Refresh()
}
func (e *CodeEditor) DragEnd() { e.dragging = false }

func (e *CodeEditor) DoubleTapped(event *fyne.PointEvent) {
	position := e.cursorOffsetFromPosition(event.Position)
	runes := []rune(e.Text())
	if position >= len(runes) {
		return
	}
	start, end := position, position+1
	if isWordRune(runes[position]) {
		for start > 0 && isWordRune(runes[start-1]) {
			start--
		}
		end = position
		for end < len(runes) && isWordRune(runes[end]) {
			end++
		}
	}
	e.cursors = []cursorState{{position: end, anchor: start}}
	e.syncPrimaryCursor()
	e.Refresh()
}

func (e *CodeEditor) cursorOffsetFromPosition(position fyne.Position) int {
	lineHeight := e.lineHeight()
	charWidth := e.charWidth()
	line := max(int(position.Y/lineHeight), 0)
	if line >= len(e.lines) {
		line = len(e.lines) - 1
	}
	column := max(int((position.X-e.lineNumberWidth())/charWidth), 0)
	return e.offsetForLineColumn(line, column)
}

func (e *CodeEditor) syncPrimaryCursor() {
	if len(e.cursors) == 0 {
		e.cursors = []cursorState{{}}
	}
	e.Buffer.SetCursor(e.cursors[0].position)
	if e.OnCursorChanged != nil {
		e.OnCursorChanged()
	}
}

func (e *CodeEditor) moveCursors(move func(int) int) {
	for index := range e.cursors {
		e.cursors[index].position = move(e.cursors[index].position)
	}
}

func (e *CodeEditor) moveCursorsLine(up bool) {
	lines := e.Buffer.Lines()
	for index := range e.cursors {
		position := e.cursors[index].position
		line, column := e.lineColumnAt(position)
		if up && line > 0 {
			line--
		} else if !up && line < len(lines)-1 {
			line++
		}
		e.cursors[index].position = e.offsetForLineColumn(line, column)
	}
}

func (e *CodeEditor) lineColumnAt(position int) (int, int) {
	e.Buffer.SetCursor(position)
	return e.Buffer.CursorLineColumn()
}

func (e *CodeEditor) replaceSelections(replacement []rune) {
	type edit struct{ index, start, end int }
	edits := make([]edit, 0, len(e.cursors))
	for index, cursor := range e.cursors {
		start, end := cursor.anchor, cursor.position
		if start > end {
			start, end = end, start
		}
		edits = append(edits, edit{index: index, start: start, end: end})
	}
	sort.SliceStable(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	for _, current := range edits {
		e.Buffer.ReplaceRange(current.start, current.end, replacement)
		newPosition := current.start + len(replacement)
		e.cursors[current.index] = cursorState{position: newPosition, anchor: newPosition}
		for index := range e.cursors {
			if index == current.index {
				continue
			}
			if e.cursors[index].position >= current.end {
				e.cursors[index].position += len(replacement) - (current.end - current.start)
			}
			if e.cursors[index].anchor >= current.end {
				e.cursors[index].anchor += len(replacement) - (current.end - current.start)
			}
		}
	}
}

func (e *CodeEditor) deleteFromCursors(forward bool) {
	for index := range e.cursors {
		if e.cursors[index].position != e.cursors[index].anchor {
			continue
		}
		if forward {
			e.cursors[index].anchor = e.cursors[index].position + 1
		} else if e.cursors[index].position > 0 {
			e.cursors[index].anchor = e.cursors[index].position - 1
		}
	}
	e.deleteSelections()
}

func (e *CodeEditor) deleteSelections() {
	e.replaceSelections(nil)
}

func (e *CodeEditor) selectedText() string {
	parts := make([]string, 0)
	for _, cursor := range e.cursors {
		start, end := cursor.anchor, cursor.position
		if start > end {
			start, end = end, start
		}
		if start != end {
			parts = append(parts, string([]rune(e.Text())[start:end]))
		}
	}
	return strings.Join(parts, "\n")
}

func (e *CodeEditor) selectAll() {
	end := len([]rune(e.Text()))
	e.cursors = []cursorState{{position: end, anchor: 0}}
	e.syncPrimaryCursor()
	e.Refresh()
}

func (e *CodeEditor) cutCurrentLines() {
	runes := []rune(e.Text())
	seen := make(map[int]bool)
	for index, cursor := range e.cursors {
		line, start, end := e.lineBoundsAt(cursor.position)
		if seen[line] {
			e.cursors[index] = cursorState{position: start, anchor: start}
			continue
		}
		seen[line] = true
		if end < len(runes) {
			end++
		} else if start > 0 {
			start--
		}
		e.cursors[index] = cursorState{position: end, anchor: start}
	}
	e.deleteSelections()
}

func (e *CodeEditor) currentLinesText() string {
	parts := make([]string, 0, len(e.cursors))
	seen := make(map[int]bool)
	for _, cursor := range e.cursors {
		line, _, _ := e.lineBoundsAt(cursor.position)
		if seen[line] {
			continue
		}
		seen[line] = true
		lines := e.Buffer.Lines()
		parts = append(parts, lines[line])
	}
	return strings.Join(parts, "\n")
}

func (e *CodeEditor) pasteLinesBelow() {
	lines := strings.Split(e.lineClipboardText, "\n")
	content := []rune(strings.Join(lines, "\n"))
	runes := []rune(e.Text())
	type insertion struct {
		position int
		text     []rune
	}
	insertions := make([]insertion, 0, len(e.cursors))
	seen := make(map[int]bool)
	for _, cursor := range e.cursors {
		line, _, lineEnd := e.lineBoundsAt(cursor.position)
		if seen[line] {
			continue
		}
		seen[line] = true
		position := lineEnd
		replacement := append([]rune{'\n'}, content...)
		if lineEnd < len(runes) {
			position++
			replacement = append(content, '\n')
		}
		insertions = append(insertions, insertion{position: position, text: replacement})
	}
	sort.Slice(insertions, func(i, j int) bool { return insertions[i].position > insertions[j].position })
	for _, current := range insertions {
		e.Buffer.ReplaceRange(current.position, current.position, current.text)
		for index := range e.cursors {
			if e.cursors[index].position >= current.position {
				e.cursors[index].position += len(current.text)
			}
			if e.cursors[index].anchor >= current.position {
				e.cursors[index].anchor += len(current.text)
			}
		}
	}
}

func (e *CodeEditor) lineBoundsAt(position int) (int, int, int) {
	runes := []rune(e.Text())
	if position > len(runes) {
		position = len(runes)
	}
	line := 0
	start := 0
	for index, character := range runes {
		if index >= position {
			break
		}
		if character == '\n' {
			line++
			start = index + 1
		}
	}
	end := start
	for end < len(runes) && runes[end] != '\n' {
		end++
	}
	return line, start, end
}

func (e *CodeEditor) selectionRectangles(cursor cursorState, lineHeight, lineNumberWidth, charWidth float32) []*canvas.Rectangle {
	start, end := cursor.anchor, cursor.position
	if start > end {
		start, end = end, start
	}
	if start == end {
		return nil
	}

	rectangles := make([]*canvas.Rectangle, 0)
	offset := 0
	for lineIndex, line := range e.Buffer.Lines() {
		lineLength := len([]rune(line))
		lineEnd := offset + lineLength
		selectionStart := max(start, offset)
		selectionEnd := min(end, lineEnd)
		if selectionStart < selectionEnd {
			localStart := selectionStart - offset
			localEnd := selectionEnd - offset
			startColumn := highlighter.VisualColumn(string([]rune(line)[:localStart]))
			endColumn := highlighter.VisualColumn(string([]rune(line)[:localEnd]))
			rectangle := canvas.NewRectangle(color.NRGBA{R: 45, G: 90, B: 150, A: 180})
			rectangle.Move(fyne.NewPos(lineNumberWidth+float32(startColumn)*charWidth, float32(lineIndex)*lineHeight))
			rectangle.Resize(fyne.NewSize(float32(endColumn-startColumn)*charWidth, lineHeight))
			rectangles = append(rectangles, rectangle)
		}
		offset = lineEnd + 1
	}
	return rectangles
}

func (e *CodeEditor) offsetForLineColumn(line, column int) int {
	lines := e.Buffer.Lines()
	offset := 0
	for i := 0; i < line && i < len(lines); i++ {
		offset += len([]rune(lines[i])) + 1
	}
	if line >= len(lines) {
		return len([]rune(e.Text()))
	}
	visual := 0
	for index, character := range []rune(lines[line]) {
		next := highlighter.VisualColumn(string(character))
		if character == '\t' {
			next = highlighter.TabWidth - visual%highlighter.TabWidth
		}
		if visual+next > column {
			return offset + index
		}
		visual += next
	}
	return offset + len([]rune(lines[line]))
}

func (e *CodeEditor) lineHeight() float32 {
	return theme.Size(theme.SizeNameText) + theme.Size(theme.SizeNameLineSpacing)
}

func (e *CodeEditor) charWidth() float32 {
	return fyne.MeasureText("M", theme.Size(theme.SizeNameText), fyne.TextStyle{Monospace: true}).Width
}

func (e *CodeEditor) lineNumberWidth() float32 {
	width := len(e.lines)
	digits := 1
	for width >= 10 {
		width /= 10
		digits++
	}
	return float32(digits+2) * e.charWidth()
}

func (e *CodeEditor) CreateRenderer() fyne.WidgetRenderer {
	return &codeEditorRenderer{editor: e}
}

type codeEditorRenderer struct {
	editor     *CodeEditor
	objects    []fyne.CanvasObject
	cursors    []*canvas.Rectangle
	selections []*canvas.Rectangle
}

func (r *codeEditorRenderer) Destroy()                     {}
func (r *codeEditorRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *codeEditorRenderer) MinSize() fyne.Size {
	maxColumns := 1
	for _, line := range r.editor.lines {
		if columns := len([]rune(line.Text)); columns > maxColumns {
			maxColumns = columns
		}
	}
	return fyne.NewSize(r.editor.lineNumberWidth()+float32(maxColumns)*r.editor.charWidth(), float32(len(r.editor.lines))*r.editor.lineHeight())
}

func (r *codeEditorRenderer) Layout(size fyne.Size) {
	charWidth := r.editor.charWidth()
	lineHeight := r.editor.lineHeight()
	lineNumberWidth := r.editor.lineNumberWidth()
	lines := r.editor.Buffer.Lines()
	for index, cursor := range r.cursors {
		if index >= len(r.editor.cursors) {
			break
		}
		position := r.editor.cursors[index].position
		r.editor.Buffer.SetCursor(position)
		line, rawColumn := r.editor.Buffer.CursorLineColumn()
		lineText := lines[line]
		column := highlighter.VisualColumn(string([]rune(lineText)[:rawColumn]))
		cursor.Resize(fyne.NewSize(2, lineHeight))
		cursor.Move(fyne.NewPos(lineNumberWidth+float32(column)*charWidth, float32(line)*lineHeight))
		if r.editor.focused {
			cursor.Show()
		} else {
			cursor.Hide()
		}
	}
	_ = size
}

func (r *codeEditorRenderer) Refresh() {
	for _, object := range r.objects {
		object.Hide()
	}
	r.objects = r.objects[:0]
	r.cursors = nil
	r.selections = nil
	lineHeight := r.editor.lineHeight()
	lineNumberWidth := r.editor.lineNumberWidth()
	charWidth := r.editor.charWidth()
	for _, cursor := range r.editor.cursors {
		for _, selection := range r.editor.selectionRectangles(cursor, lineHeight, lineNumberWidth, charWidth) {
			r.selections = append(r.selections, selection)
			r.objects = append(r.objects, selection)
		}
	}
	for lineIndex, line := range r.editor.lines {
		y := float32(lineIndex) * lineHeight
		for _, span := range line.Spans {
			if background := spanBackground(span.Kind); background != nil {
				width := float32(len([]rune(line.Text))) * charWidth
				if available := r.editor.Size().Width - lineNumberWidth; available > width {
					width = available
				}
				backgroundRectangle := canvas.NewRectangle(background)
				backgroundRectangle.Move(fyne.NewPos(lineNumberWidth, y))
				backgroundRectangle.Resize(fyne.NewSize(width, lineHeight))
				r.objects = append(r.objects, backgroundRectangle)
			}
		}
		lineNumber := canvas.NewText(fmt.Sprintf("%d", lineIndex+1), theme.Color(theme.ColorNameDisabled))
		lineNumber.TextStyle.Monospace = true
		lineNumber.Move(fyne.NewPos(0, y))
		r.objects = append(r.objects, lineNumber)
		cursor := 0
		for _, span := range line.Spans {
			if spanBackground(span.Kind) != nil {
				continue
			}
			if span.Start > len([]rune(line.Text)) {
				continue
			}
			if span.Start > cursor {
				r.objects = append(r.objects, newTextSpan(string([]rune(line.Text)[cursor:span.Start]), theme.Color(theme.ColorNameForeground), lineNumberWidth+float32(cursor)*charWidth, y))
			}
			end := min(span.End, len([]rune(line.Text)))
			r.objects = append(r.objects, newTextSpan(string([]rune(line.Text)[span.Start:end]), spanColor(span.Kind), lineNumberWidth+float32(span.Start)*charWidth, y))
			cursor = end
		}
		if cursor < len([]rune(line.Text)) {
			r.objects = append(r.objects, newTextSpan(string([]rune(line.Text)[cursor:]), theme.Color(theme.ColorNameForeground), lineNumberWidth+float32(cursor)*charWidth, y))
		}
	}
	for range r.editor.cursors {
		cursor := canvas.NewRectangle(cursorColor)
		r.cursors = append(r.cursors, cursor)
		r.objects = append(r.objects, cursor)
	}
	r.Layout(r.editor.Size())
}

func spanColor(kind highlighter.HighlightKind) color.Color {
	switch kind {
	case highlighter.HighlightComment:
		return commentColor
	case highlighter.HighlightString:
		return stringColor
	case highlighter.HighlightNumber:
		return numberColor
	default:
		return keywordColor
	}
}

func spanBackground(kind highlighter.HighlightKind) color.Color {
	switch kind {
	case highlighter.HighlightDiffAdded:
		return color.NRGBA{R: 70, G: 160, B: 80, A: 90}
	case highlighter.HighlightDiffRemoved:
		return color.NRGBA{R: 210, G: 70, B: 70, A: 100}
	default:
		return nil
	}
}

func newTextSpan(text string, foreground color.Color, x, y float32) *canvas.Text {
	span := canvas.NewText(text, foreground)
	span.TextStyle.Monospace = true
	span.Move(fyne.NewPos(x, y))
	return span
}
