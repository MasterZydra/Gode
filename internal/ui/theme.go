package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type folderTheme struct {
	base fyne.Theme
}

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

func NewFolderTheme() fyne.Theme {
	return folderTheme{base: theme.DefaultTheme()}
}
