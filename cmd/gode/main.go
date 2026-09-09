package main

import (
	"gode/internal/editor"
	"gode/internal/explorer"
	"gode/internal/ui"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(ui.NewFolderTheme())
	myWindow := myApp.NewWindow("Table Widget")
	myWindow.Resize(fyne.NewSize(1440, 801))

	fileExplorer := explorer.NewExplorer("")
	if len(os.Args) > 1 {
		rootPath := os.Args[1]
		if fileInfo, err := os.Stat(rootPath); err == nil && fileInfo.IsDir() {
			fileExplorer = explorer.NewExplorer(rootPath)
			if err := fileExplorer.Update(); err != nil {
				fileExplorer = explorer.NewExplorer("")
			}
		}
	}

	textEditor := editor.New()
	updateWindowTitle := func() {
		myWindow.SetTitle(textEditor.Title())
	}
	textEditor.Widget.OnChanged = func(string) { updateWindowTitle() }
	updateWindowTitle()
	saveShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	myWindow.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File", &fyne.MenuItem{
		Label: "Save",
		Action: func() {
			if err := textEditor.Save(); err != nil {
				dialog.ShowError(err, myWindow)
			}
			updateWindowTitle()
		},
		Shortcut: saveShortcut,
	})))
	tree := ui.NewTree(fileExplorer, func(node *explorer.Node) {
		if node.IsDir {
			textEditor.Clear()
			updateWindowTitle()
			return
		}

		if err := textEditor.Load(node.Path); err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		updateWindowTitle()
	})
	textEditorContainer := container.New(layout.NewCustomPaddedLayout(10, 10, 10, 10), textEditor.Widget)
	myWindow.SetContent(container.NewBorder(nil, nil, tree, nil, textEditorContainer))
	myWindow.ShowAndRun()
}
