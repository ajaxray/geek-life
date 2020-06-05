package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func prepareProjectDetail() {
	deleteBtn := makeButton("Delete Project", deleteCurrentProject)
	clearBtn := makeButton("Clear Completed Tasks", clearCompletedTasks)

	deleteBtn.SetBackgroundColor(tcell.ColorRed)
	projectDetailPane = tview.NewFlex().SetDirection(tview.FlexRow).
		// AddItem(activeProjectName, 1, 1, false).
		// AddItem(makeHorizontalLine(tcell.RuneS3, tcell.ColorGray), 1, 1, false).
		AddItem(deleteBtn, 3, 1, false).
		AddItem(blankCell, 1, 1, false).
		AddItem(clearBtn, 3, 1, false).
		AddItem(blankCell, 0, 1, false)

	projectDetailPane.SetBorder(true).SetTitle("[::u]A[::-]ctions")
}

func deleteCurrentProject() {
	if currentProject != nil && projectRepo.Delete(currentProject) == nil {
		for i, _ := range tasks {
			taskRepo.Delete(&tasks[i])
		}

		showMessage("Removed Project: " + currentProject.Title)
		removeThirdCol()
		taskList.Clear()
		projectList.Clear()

		loadProjectList()
	}
}

func clearCompletedTasks() {
	count := 0
	for i, task := range tasks {
		if task.Completed && taskRepo.Delete(&tasks[i]) == nil {
			taskList.RemoveItem(i)
			count++
		}
	}
	showMessage(fmt.Sprintf("[yellow]%d tasks cleared!", count))
}
