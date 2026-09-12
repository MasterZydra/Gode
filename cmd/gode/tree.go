package main

import (
	"gode/internal/explorer"
	"gode/internal/ui"
	"path/filepath"

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
	gitView = ui.NewGitView(fileExplorer.RootDir(), func(path string) {
		withSavedChanges(myWindow, func() {
			if err := textEditor.Load(filepath.Join(fileExplorer.RootDir(), path)); err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			editorScroll.ScrollToTop()
			myWindow.SetTitle(textEditor.Title())
		})
	}, func(err error) {
		dialog.ShowError(err, myWindow)
	})
	statusBar = ui.NewStatusBar(fileExplorer.RootDir())
	gitView.SetOnRefresh(statusBar.Refresh)
	textEditorContainer := container.New(layout.NewCustomPaddedLayout(10, 10, 10, 10), editorScroll)
	sidebar := container.NewAppTabs(
		container.NewTabItem("Explorer", fileTree),
		container.NewTabItem("Git", gitView.CanvasObject()),
	)
	sidebar.SetTabLocation(container.TabLocationLeading)
	content := container.NewHSplit(sidebar, textEditorContainer)
	content.SetOffset(0.25)
	myWindow.SetContent(container.NewBorder(nil, statusBar.CanvasObject(), nil, nil, content))
}
