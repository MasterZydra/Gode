package main

import (
	"gode/internal/explorer"
	"gode/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
)

func setTree(myWindow fyne.Window) {
	fileTree = ui.NewTree(fileExplorer, func(node *explorer.Node) {
		if !node.IsDir && node.Path == textEditor.SelectedPath {
			return
		}
		withSavedChanges(myWindow, func() {
			if node.IsDir {
				textEditor.Clear()
				editorScroll.ScrollToTop()
				myWindow.SetTitle(textEditor.Title())
				return
			}
			if err := textEditor.Load(node.Path); err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			editorScroll.ScrollToTop()
			myWindow.SetTitle(textEditor.Title())
		})
	})
	textEditorContainer := container.New(layout.NewCustomPaddedLayout(10, 10, 10, 10), editorScroll)
	myWindow.SetContent(container.NewBorder(nil, nil, fileTree, nil, textEditorContainer))
}
