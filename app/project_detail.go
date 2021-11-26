package main

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
)

// ProjectDetailPane Displays relevant actions of current project
type ProjectDetailPane struct {
	*tview.Flex
	project *model.Project
}

// NewProjectDetailPane Initializes ProjectDetailPane
func NewProjectDetailPane() *ProjectDetailPane {
	pane := ProjectDetailPane{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
	}
	deleteBtn := makeButton("[::u]D[::-]elete Project", projectPane.RemoveActivateProject)
	clearBtn := makeButton("[::u]C[::-]lear Completed Tasks", taskPane.ClearCompletedTasks)

	deleteBtn.SetBackgroundColor(tcell.ColorRed)
	pane.
		AddItem(deleteBtn, 3, 1, false).
		AddItem(blankCell, 1, 1, false).
		AddItem(clearBtn, 3, 1, false).
		AddItem(blankCell, 0, 1, false)

	pane.SetBorder(true).SetTitle("[::u]A[::-]ctions")

	return &pane
}

// SetProject Sets the active Project
func (pd *ProjectDetailPane) SetProject(project *model.Project) {
	pd.project = project
	pd.SetTitle("[::b]" + pd.project.Title)
}

func (pd *ProjectDetailPane) isShowing() bool {
	return taskPane.activeTask == nil && projectPane.activeProject != nil
}

func (pd *ProjectDetailPane) handleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch unicode.ToLower(event.Rune()) {
	case 'd':
		projectPane.RemoveActivateProject()
		return nil
	case 'c':
		taskPane.ClearCompletedTasks()
		return nil
	}

	return event
}
