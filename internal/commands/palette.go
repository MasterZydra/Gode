package commands

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Palette struct {
	registry *Registry
	window   fyne.Window
	entry    *widget.Entry
	list     *widget.List
	results  []Command
	dialog   *dialog.CustomDialog
}

func NewPalette(registry *Registry, window fyne.Window) *Palette {
	palette := &Palette{registry: registry, window: window}
	palette.entry = widget.NewEntry()
	palette.entry.SetPlaceHolder("Search commands")
	palette.list = widget.NewList(
		func() int { return len(palette.results) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(palette.results[id].Name)
		},
	)
	palette.entry.OnChanged = func(query string) {
		palette.results = registry.Search(query)
		palette.list.Refresh()
	}
	palette.entry.OnSubmitted = func(string) {
		palette.executeFirst()
	}
	palette.list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(palette.results) {
			palette.execute(palette.results[id])
		}
	}
	palette.dialog = dialog.NewCustomWithoutButtons(
		"Command Palette",
		container.NewVBox(palette.entry, palette.list),
		window,
	)
	return palette
}

func (p *Palette) Show() {
	p.entry.SetText("")
	p.results = p.registry.Search("")
	p.list.Refresh()
	p.dialog.Show()
	if canvas := fyne.CurrentApp().Driver().CanvasForObject(p.entry); canvas != nil {
		canvas.Focus(p.entry)
	}
}

func (p *Palette) executeFirst() {
	if len(p.results) > 0 {
		p.execute(p.results[0])
	}
}

func (p *Palette) execute(command Command) {
	p.dialog.Hide()
	if command.Action != nil {
		command.Action()
	}
}
