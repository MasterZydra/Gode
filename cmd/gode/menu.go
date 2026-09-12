package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
)

func setMainMenu(myWindow fyne.Window) {
	commandPalette := registerEditorCommands(myWindow)

	fileMenu := fyne.NewMenu("File",
		&fyne.MenuItem{Label: "Open Folder", Action: func() {
			dialog.ShowFolderOpen(func(selected fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				if selected == nil {
					return
				}
				withSavedChanges(myWindow, func() {
					if err := fileExplorer.SetRoot(selected.Path()); err != nil {
						dialog.ShowError(err, myWindow)
						return
					}
					fileTree.Refresh()
					gitView.SetRoot(selected.Path())
					textEditor.Clear()
				})
			}, myWindow)
		}, Shortcut: &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}},
		&fyne.MenuItem{Label: "Save", Action: func() {
			if err := textEditor.Save(); err != nil {
				dialog.ShowError(err, myWindow)
			}
		}, Shortcut: &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}},
	)

	toolsMenu := fyne.NewMenu("Tools",
		&fyne.MenuItem{Label: "Command Palette", Action: commandPalette.Show,
			Shortcut: &desktop.CustomShortcut{KeyName: fyne.KeyP, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}},
	)

	myWindow.SetMainMenu(fyne.NewMainMenu(fileMenu, toolsMenu))
}
