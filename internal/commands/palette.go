package commands

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Palette struct {
	registry *Registry
	window   fyne.Window
	entry    *paletteEntry
	list     *widget.List
	results  []Command
	dialog   *dialog.CustomDialog
}

type paletteEntry struct {
	*widget.Entry
	onEscape func()
}

func newPaletteEntry(onEscape func()) *paletteEntry {
	entry := &paletteEntry{
		Entry:    widget.NewEntry(),
		onEscape: onEscape,
	}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *paletteEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyEscape {
		e.onEscape()
		return
	}
	e.Entry.TypedKey(key)
}

func (e *paletteEntry) MouseDown(event *desktop.MouseEvent) {
	if canvas := fyne.CurrentApp().Driver().CanvasForObject(e); canvas != nil {
		canvas.Focus(e)
	}
	e.Entry.MouseDown(event)
}

func NewPalette(registry *Registry, window fyne.Window) *Palette {
	palette := &Palette{registry: registry, window: window}
	palette.entry = newPaletteEntry(func() { palette.dialog.Hide() })
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
	p.list.UnselectAll()
	p.list.Refresh()
	p.dialog.Show()
	p.window.Canvas().Focus(p.entry)
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
