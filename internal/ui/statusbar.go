package ui

import (
	"fmt"
	"gode/internal/git"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	repository *git.Repository
	label      *widget.Label
	root       fyne.CanvasObject
}

func NewStatusBar(rootPath string) *StatusBar {
	bar := &StatusBar{
		label: widget.NewLabel(""),
	}
	bar.root = container.NewHBox(bar.label)
	bar.SetRoot(rootPath)
	return bar
}

func (b *StatusBar) SetRoot(rootPath string) {
	b.repository = git.NewRepository(rootPath)
	b.Refresh()
}

func (b *StatusBar) Refresh() {
	info, err := b.repository.Info()
	if err != nil {
		b.label.SetText(fmt.Sprintf("Git: %v", err))
		return
	}
	b.label.SetText(fmt.Sprintf("%s · %d ahead · %d behind", info.Branch, info.Ahead, info.Behind))
}

func (b *StatusBar) CanvasObject() fyne.CanvasObject {
	return b.root
}
