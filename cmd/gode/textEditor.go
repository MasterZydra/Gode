package main

import (
	"gode/internal/commands"
	"gode/internal/editor"
	"gode/internal/explorer"
	"gode/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	fileExplorer *explorer.Explorer
	textEditor   *editor.Editor
	editorScroll *container.Scroll
	fileTree     *widget.Tree
	gitView      *ui.GitView
	statusBar    *ui.StatusBar
)

func setTextEditor(myWindow fyne.Window, currentExplorer *explorer.Explorer) {
	fileExplorer = currentExplorer
	textEditor = editor.New()
	editorScroll = container.NewScroll(textEditor.Widget)
	textEditor.OnStateChanged = func() { myWindow.SetTitle(textEditor.Title()) }
	myWindow.SetTitle(textEditor.Title())
	textEditor.Widget.OnCursorChanged = func() {
		if editorScroll.Size().Height == 0 {
			return
		}
		cursorPosition, cursorSize := textEditor.Widget.CursorBounds()
		offset := editorScroll.Offset
		if cursorPosition.Y < offset.Y {
			offset.Y = cursorPosition.Y
		} else if cursorPosition.Y+cursorSize.Height > offset.Y+editorScroll.Size().Height {
			offset.Y = cursorPosition.Y + cursorSize.Height - editorScroll.Size().Height
		}
		editorScroll.ScrollToOffset(offset)
	}
	myWindow.SetCloseIntercept(func() {
		closeWindow := func() {
			myWindow.SetCloseIntercept(nil)
			myWindow.Close()
		}
		if !textEditor.Dirty {
			closeWindow()
			return
		}
		showUnsavedChanges(myWindow, "Save changes to "+textEditor.FileName()+" before closing?", closeWindow)
	})
}

func showUnsavedChanges(myWindow fyne.Window, message string, next func()) {
	prompt := dialog.NewCustomWithoutButtons("Unsaved changes", widget.NewLabel(message), myWindow)
	discard := widget.NewButton("Discard changes", func() { prompt.Hide(); next() })
	cancel := widget.NewButton("Cancel", prompt.Hide)
	save := widget.NewButton("Save", func() {
		if err := textEditor.Save(); err != nil {
			prompt.Hide()
			dialog.ShowError(err, myWindow)
			return
		}
		prompt.Hide()
		next()
	})
	save.Importance = widget.HighImportance
	prompt.SetButtons([]fyne.CanvasObject{discard, cancel, save})
	prompt.Show()
}

func withSavedChanges(myWindow fyne.Window, next func()) {
	if !textEditor.Dirty {
		next()
		return
	}
	showUnsavedChanges(myWindow, "Save changes to "+textEditor.FileName()+" before continuing?", next)
}

func registerEditorCommands(myWindow fyne.Window) *commands.Palette {
	commandRegistry := commands.NewRegistry()
	commandPalette := commands.NewPalette(commandRegistry, myWindow)
	commandRegistry.Register(commands.Command{
		Name: "Format Document",
		Action: func() {
			if err := textEditor.Format(); err != nil {
				dialog.ShowError(err, myWindow)
			}
		},
	})
	return commandPalette
}
