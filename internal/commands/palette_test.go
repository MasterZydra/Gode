package commands

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestPaletteEntryEscape(t *testing.T) {
	dismissed := false
	entry := newPaletteEntry(func() { dismissed = true })

	entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})

	if !dismissed {
		t.Fatal("Escape did not dismiss the palette")
	}
}

func TestPaletteEntryMouseDownFocusesEntry(t *testing.T) {
	test.NewTempApp(t)
	entry := newPaletteEntry(nil)
	window := test.NewTempWindow(t, entry)

	entry.MouseDown(&desktop.MouseEvent{
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 1)},
	})

	if window.Canvas().Focused() != entry {
		t.Fatal("clicking the palette entry did not focus it")
	}
}

func TestPaletteShowClearsPreviousSelection(t *testing.T) {
	test.NewTempApp(t)
	window := test.NewTempWindow(t, nil)
	registry := NewRegistry()
	registry.Register(Command{Name: "Format Document"})
	palette := NewPalette(registry, window)
	selected := 0
	palette.list.OnSelected = func(widget.ListItemID) { selected++ }

	palette.Show()
	palette.list.Select(0)
	palette.dialog.Hide()
	palette.Show()
	palette.list.Select(0)

	if selected != 2 {
		t.Fatalf("command selected %d times, want 2", selected)
	}
}
