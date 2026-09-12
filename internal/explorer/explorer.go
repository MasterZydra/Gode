package explorer

import (
	"os"
	"sort"
)

type Explorer struct {
	rootDir string
	Nodes   []*Node
}

func NewExplorer(rootDir string) *Explorer {
	return &Explorer{
		rootDir: rootDir,
	}
}

func (explorer *Explorer) RootDir() string {
	return explorer.rootDir
}

func (explorer *Explorer) Update() error {
	nodes, err := readNodes(explorer.rootDir)
	if err != nil {
		return err
	}

	explorer.Nodes = nodes
	return nil
}

func (explorer *Explorer) SetRoot(rootDir string) error {
	explorer.rootDir = rootDir
	return explorer.Update()
}

func readNodes(dir string) ([]*Node, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	nodes := make([]*Node, 0, len(entries))
	for _, entry := range entries {
		node := NodeFromOsDirEntry(dir, entry)
		if node.IsDir {
			node.Children, err = readNodes(node.Path)
			if err != nil {
				return nil, err
			}
		}
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].IsDir != nodes[j].IsDir {
			return nodes[i].IsDir
		}
		return nodes[i].Name < nodes[j].Name
	})

	return nodes, nil
}
