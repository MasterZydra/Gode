package main

import (
	"gode/internal/explorer"
	"gode/internal/ui/theme"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.gode.editor")
	myApp.Settings().SetTheme(theme.NewFolderTheme())
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

	setTextEditor(myWindow, fileExplorer)
	setTree(myWindow)
	setMainMenu(myWindow)
	myWindow.ShowAndRun()
}
