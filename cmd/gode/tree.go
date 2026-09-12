package main

import (
	"gode/internal/explorer"
	"gode/internal/highlighter"
	"gode/internal/ui/git"
	"gode/internal/ui/statusbar"
	"gode/internal/ui/tree"
	"gode/internal/widgets/codeeditor"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func setTree(myWindow fyne.Window) {
	fileTree = tree.NewTree(fileExplorer, func(node *explorer.Node) {
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
	gitView = git.NewGitView(fileExplorer.RootDir(), func(path string) {
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
	gitView.SetOnDiff(func(path, output string) {
		diffView := fyne.CurrentApp().NewWindow("Diff: " + path)
		diffText := codeeditor.NewCodeEditor()
		diffText.SetHighlighter(highlighter.NewDiffHighlighter(highlighter.HighlighterForFile(path)))
		diffText.SetReadOnly(true)
		diffText.SetText(output)
		diffView.SetContent(container.NewBorder(nil, widget.NewButton("Close", diffView.Close), nil, nil, container.NewScroll(diffText)))
		diffView.Resize(fyne.NewSize(900, 600))
		diffView.Show()
	})
	statusBar = statusbar.NewStatusBar(fileExplorer.RootDir())
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
