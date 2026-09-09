package tabs

type Tab struct {
	ID       string
	FilePath string
	Editor   *editor.Editor
}
