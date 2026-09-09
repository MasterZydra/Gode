package explorer

import (
	"os"
	"path"
)

var nodeId int64 = 0

func nextNodeId() int64 {
	nodeId++
	return nodeId
}

type Node struct {
	Id       int64
	Name     string
	Path     string
	IsDir    bool
	Children []*Node
}

func NodeFromOsDirEntry(rootDir string, dirEntry os.DirEntry) *Node {
	return &Node{
		Id:       nextNodeId(),
		Name:     dirEntry.Name(),
		Path:     path.Join(rootDir, dirEntry.Name()),
		IsDir:    dirEntry.IsDir(),
		Children: []*Node{},
	}
}
