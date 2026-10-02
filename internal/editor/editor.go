package editor

import (
	"fmt"
	"gode/internal/formatter"
	"gode/internal/highlighter"
	"gode/internal/widgets/codeeditor"
	"os"
	"path/filepath"
	"strings"
)

const MaxFileSize int64 = 5 * 1024 * 1024

type Editor struct {
	Widget         *codeeditor.CodeEditor
	SelectedPath   string
	Dirty          bool
	OnStateChanged func()
	loading        bool
	lineEnding     string
}

func New() *Editor {
	codeEditor := codeeditor.NewCodeEditor()
	editor := &Editor{Widget: codeEditor, lineEnding: "\n"}
	codeEditor.OnChanged = func(string) {
		if !editor.loading {
			editor.Dirty = true
		}
		if editor.OnStateChanged != nil {
			editor.OnStateChanged()
		}
	}
	return editor
}

func (e *Editor) Load(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fileInfo.Size() > MaxFileSize {
		return fmt.Errorf("file is too large to display (maximum %d MiB)", MaxFileSize/(1024*1024))
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	text, lineEnding := normalizeLineEndings(contents)
	e.loading = true
	e.Widget.SetHighlighter(highlighter.HighlighterForFile(path))
	e.Widget.SetText(text)
	e.loading = false
	e.SelectedPath = path
	e.lineEnding = lineEnding
	e.Dirty = false
	e.notifyStateChanged()
	return nil
}

func (e *Editor) Save() error {
	if e.SelectedPath == "" {
		return nil
	}

	fileInfo, err := os.Stat(e.SelectedPath)
	if err != nil {
		return err
	}
	text := e.Widget.Text()
	if e.lineEnding != "\n" {
		text = strings.ReplaceAll(text, "\n", e.lineEnding)
	}
	if err := os.WriteFile(e.SelectedPath, []byte(text), fileInfo.Mode().Perm()); err != nil {
		return err
	}

	e.Dirty = false
	e.notifyStateChanged()
	return nil
}

func (e *Editor) Format() error {
	if e.SelectedPath == "" {
		return fmt.Errorf("no file selected")
	}

	formatted, err := formatter.Format(e.SelectedPath, e.Widget.Text())
	if err != nil {
		return err
	}
	if formatted == e.Widget.Text() {
		return nil
	}
	e.Widget.SetText(formatted)
	e.Dirty = true
	e.notifyStateChanged()
	return nil
}

func (e *Editor) Clear() {
	e.SelectedPath = ""
	e.Dirty = false
	e.lineEnding = "\n"
	e.loading = true
	e.Widget.SetHighlighter(nil)
	e.Widget.SetText("")
	e.loading = false
	e.notifyStateChanged()
}

func (e *Editor) notifyStateChanged() {
	if e.OnStateChanged != nil {
		e.OnStateChanged()
	}
}

func (e *Editor) FileName() string {
	return filepath.Base(e.SelectedPath)
}

func (e *Editor) Title() string {
	if e.SelectedPath == "" {
		return "Gode"
	}
	if e.Dirty {
		return "• " + e.FileName()
	}
	return e.FileName()
}

func normalizeLineEndings(contents []byte) (string, string) {
	var text strings.Builder
	text.Grow(len(contents))
	lineEnding := "\n"
	foundLineEnding := false

	for i := 0; i < len(contents); i++ {
		switch contents[i] {
		case '\r':
			lineEndingForCharacter := "\r"
			if i+1 < len(contents) && contents[i+1] == '\n' {
				lineEndingForCharacter = "\r\n"
				i++
			}
			if !foundLineEnding {
				lineEnding = lineEndingForCharacter
				foundLineEnding = true
			}
			text.WriteByte('\n')
		case '\n':
			if !foundLineEnding {
				lineEnding = "\n"
				foundLineEnding = true
			}
			text.WriteByte('\n')
		default:
			text.WriteByte(contents[i])
		}
	}

	return text.String(), lineEnding
}
