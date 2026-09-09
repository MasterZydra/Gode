package ui

import (
	"gode/internal/explorer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewTree(fileExplorer *explorer.Explorer, onFileSelected func(*explorer.Node)) *widget.Tree {
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
		func(id widget.TreeNodeID, _ bool, o fyne.CanvasObject) {
			text := string(id)
			if node := findNode(id); node != nil {
				text = node.Name
			}
			o.(*widget.Label).SetText(text)
		})
	tree.OnSelected = func(id widget.TreeNodeID) {
		node := findNode(id)
		if node == nil {
			return
		}
		if node.IsDir {
			tree.ToggleBranch(id)
		}
		if onFileSelected != nil {
			onFileSelected(node)
		}
	}

	return tree
}
