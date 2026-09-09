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
	myApp := app.NewWithID("com.gode.editor")
	myApp.Settings().SetTheme(ui.NewFolderTheme())
	myWindow := myApp.NewWindow("Gode")
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
	textEditor.OnStateChanged = updateWindowTitle
	updateWindowTitle()
	withSavedChanges := func(next func()) {
		if !textEditor.Dirty {
			next()
			return
		}
		dialog.ShowConfirm("Save changes?", "Save changes to "+textEditor.FileName()+" before continuing?", func(save bool) {
			if !save {
				return
			}
			if err := textEditor.Save(); err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			next()
		}, myWindow)
	}
	tree := ui.NewTree(fileExplorer, func(node *explorer.Node) {
		if !node.IsDir && node.Path == textEditor.SelectedPath {
			return
		}
		withSavedChanges(func() {
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
	})
	myWindow.SetCloseIntercept(func() {
		closeWindow := func() {
			myWindow.SetCloseIntercept(nil)
			myWindow.Close()
		}
		if !textEditor.Dirty {
			closeWindow()
			return
		}
		dialog.ShowConfirm("Save changes?", "Save changes to "+textEditor.FileName()+" before closing?", func(save bool) {
			if !save {
				return
			}
			if err := textEditor.Save(); err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			closeWindow()
		}, myWindow)
	})
	myWindow.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		&fyne.MenuItem{
			Label: "Open Folder",
			Action: func() {
				dialog.ShowFolderOpen(func(selected fyne.ListableURI, err error) {
					if err != nil {
						dialog.ShowError(err, myWindow)
						return
					}
					if selected == nil {
						return
					}
					withSavedChanges(func() {
						if err := fileExplorer.SetRoot(selected.Path()); err != nil {
							dialog.ShowError(err, myWindow)
							return
						}
						tree.Refresh()
						textEditor.Clear()
						updateWindowTitle()
					})
				}, myWindow)
			},
			Shortcut: &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl},
		},
		&fyne.MenuItem{
			Label: "Save",
			Action: func() {
				if err := textEditor.Save(); err != nil {
					dialog.ShowError(err, myWindow)
				}
				updateWindowTitle()
			},
			Shortcut: &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl},
		})))
	textEditorContainer := container.New(layout.NewCustomPaddedLayout(10, 10, 10, 10), container.NewScroll(textEditor.Widget))
	myWindow.SetContent(container.NewBorder(nil, nil, tree, nil, textEditorContainer))
	myWindow.ShowAndRun()
}
