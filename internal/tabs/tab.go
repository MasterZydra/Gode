package tabs

import "gode/internal/editor"

type Tab struct {
	ID       string
	FilePath string
	Editor   *editor.Editor
}
