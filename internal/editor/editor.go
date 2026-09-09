package editor

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const MaxFileSize int64 = 5 * 1024 * 1024

type Editor struct {
	Widget       *widget.Entry
	SelectedPath string
	Dirty        bool
	loading      bool
}

func New() *Editor {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = fyne.ScrollBoth
	entry.TextStyle.Monospace = true

	editor := &Editor{Widget: entry}
	entry.OnChanged = func(string) {
		if !editor.loading {
			editor.Dirty = true
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

	e.loading = true
	e.Widget.SetText(string(contents))
	e.loading = false
	e.SelectedPath = path
	e.Dirty = false
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
	if err := os.WriteFile(e.SelectedPath, []byte(e.Widget.Text), fileInfo.Mode().Perm()); err != nil {
		return err
	}

	e.Dirty = false
	return nil
}

func (e *Editor) Clear() {
	e.SelectedPath = ""
	e.Dirty = false
	e.loading = true
	e.Widget.SetText("")
	e.loading = false
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
