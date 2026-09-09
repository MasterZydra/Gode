package main

import (
	"fmt"
	"gode/internal/explorer"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type folderTheme struct {
	base fyne.Theme
}

const maxEditorFileSize int64 = 5 * 1024 * 1024

func (t folderTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.base.Color(name, variant)
}

func (t folderTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t folderTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	switch name {
	case theme.IconNameNavigateNext:
		return theme.FolderIcon()
	case theme.IconNameMoveDown:
		return theme.FolderOpenIcon()
	default:
		return t.base.Icon(name)
	}
}

func (t folderTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(folderTheme{base: theme.DefaultTheme()})
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

	findNode := func(id widget.TreeNodeID) *explorer.Node {
		var find func(nodes []*explorer.Node) *explorer.Node
		find = func(nodes []*explorer.Node) *explorer.Node {
			for _, node := range nodes {
				if node.Path == string(id) {
					return node
				}
				if found := find(node.Children); found != nil {
					return found
				}
			}
			return nil
		}
		return find(fileExplorer.Nodes)
	}

	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			var nodes []*explorer.Node
			if id == "" {
				nodes = fileExplorer.Nodes
			} else if node := findNode(id); node != nil {
				nodes = node.Children
			}

			children := make([]widget.TreeNodeID, 0, len(nodes))
			for _, node := range nodes {
				children = append(children, widget.TreeNodeID(node.Path))
			}
			return children
		},
		func(id widget.TreeNodeID) bool {
			return id == "" || func() bool {
				node := findNode(id)
				return node != nil && node.IsDir
			}()
		},
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("Branch template")
			}
			return widget.NewLabel("Leaf template")
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			text := string(id)
			if node := findNode(id); node != nil {
				text = node.Name
			}
			// if branch {
			// 	text += " (branch)"
			// }
			o.(*widget.Label).SetText(text)
		})

	textEditor := widget.NewMultiLineEntry()
	textEditor.Wrapping = fyne.TextWrapOff
	textEditor.Scroll = fyne.ScrollBoth
	textEditor.TextStyle.Monospace = true
	selectedFileName := ""
	dirty := false
	loadingFile := false
	updateWindowTitle := func() {
		if selectedFileName == "" {
			myWindow.SetTitle("Gode")
			return
		}
		if dirty {
			myWindow.SetTitle("• " + selectedFileName)
			return
		}
		myWindow.SetTitle(selectedFileName)
	}
	textEditor.OnChanged = func(string) {
		if loadingFile {
			return
		}
		dirty = true
		updateWindowTitle()
	}
	tree.OnSelected = func(id widget.TreeNodeID) {
		node := findNode(id)
		if node == nil {
			return
		}
		if node.IsDir {
			tree.ToggleBranch(id)
			selectedFileName = ""
			dirty = false
			updateWindowTitle()
			return
		}

		fileInfo, err := os.Stat(node.Path)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if fileInfo.Size() > maxEditorFileSize {
			dialog.ShowError(fmt.Errorf("file is too large to display (maximum %d MiB)", maxEditorFileSize/(1024*1024)), myWindow)
			return
		}

		contents, err := os.ReadFile(node.Path)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		loadingFile = true
		textEditor.SetText(string(contents))
		loadingFile = false
		selectedFileName = node.Name
		dirty = false
		updateWindowTitle()
	}
	textEditorContainer := container.New(layout.NewCustomPaddedLayout(10, 10, 10, 10), textEditor)
	myWindow.SetContent(container.NewBorder(nil, nil, tree, nil, textEditorContainer))
	myWindow.ShowAndRun()
}
