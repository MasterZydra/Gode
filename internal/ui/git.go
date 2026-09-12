package ui

import (
	"fmt"
	"gode/internal/git"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GitView struct {
	repository *git.Repository
	onOpen     func(string)
	onError    func(error)

	root           *fyne.Container
	message        *widget.Entry
	commit         *widget.Button
	stagedSection  *fyne.Container
	changesSection *fyne.Container
	statusLabel    *widget.Label
}

func NewGitView(rootPath string, onOpen func(string), onError func(error)) *GitView {
	view := &GitView{
		onOpen:  onOpen,
		onError: onError,
	}
	view.message = widget.NewEntry()
	view.message.MultiLine = true
	view.message.SetMinRowsVisible(3)
	view.message.SetPlaceHolder("Commit message")
	view.message.OnChanged = func(message string) {
		view.commit.Enable()
		if strings.TrimSpace(message) == "" {
			view.commit.Disable()
		}
	}
	view.commit = widget.NewButton("Commit", view.commitChanges)
	view.commit.Disable()
	view.stagedSection = container.NewVBox()
	view.changesSection = container.NewVBox()
	view.statusLabel = widget.NewLabel("")
	view.root = container.NewVBox(
		view.message,
		view.commit,
		widget.NewSeparator(),
		view.statusLabel,
		view.stagedSection,
		view.changesSection,
	)
	view.SetRoot(rootPath)
	return view
}

func (v *GitView) SetRoot(rootPath string) {
	v.repository = git.NewRepository(rootPath)
	v.Refresh()
}

func (v *GitView) CanvasObject() fyne.CanvasObject {
	return container.NewScroll(v.root)
}

func (v *GitView) Refresh() {
	v.stagedSection.Objects = nil
	v.changesSection.Objects = []fyne.CanvasObject{widget.NewLabel("Changes")}
	v.statusLabel.SetText("")

	status, err := v.repository.Status()
	if err != nil {
		v.statusLabel.SetText(fmt.Sprintf("Git: %v", err))
		v.refreshContainers()
		return
	}

	if len(status.Staged) > 0 {
		v.stagedSection.Objects = append(v.stagedSection.Objects, widget.NewLabel("Staged Changes"))
		for _, change := range status.Staged {
			v.stagedSection.Objects = append(v.stagedSection.Objects, v.changeRow(change, "-", v.repository.Unstage))
		}
	}
	if len(status.Unstaged) == 0 {
		v.changesSection.Objects = append(v.changesSection.Objects, widget.NewLabel("No changes"))
	} else {
		for _, change := range status.Unstaged {
			v.changesSection.Objects = append(v.changesSection.Objects, v.changeRow(change, "+", v.repository.Stage))
		}
	}
	v.refreshContainers()
}

func (v *GitView) changeRow(change git.Change, actionLabel string, action func(string) error) fyne.CanvasObject {
	open := widget.NewButton(change.Path, func() {
		if v.onOpen != nil {
			v.onOpen(change.Path)
		}
	})
	actionButton := widget.NewButton(actionLabel, func() {
		if err := action(change.Path); err != nil {
			v.reportError(err)
			return
		}
		v.Refresh()
	})
	return container.NewBorder(nil, nil, nil, actionButton, open)
}

func (v *GitView) commitChanges() {
	if err := v.repository.Commit(v.message.Text); err != nil {
		v.reportError(err)
		return
	}
	v.message.SetText("")
	v.Refresh()
}

func (v *GitView) refreshContainers() {
	v.stagedSection.Refresh()
	v.changesSection.Refresh()
	v.root.Refresh()
}

func (v *GitView) reportError(err error) {
	if v.onError != nil {
		v.onError(err)
	}
}
