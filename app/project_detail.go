package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func prepareProjectDetail() {
	deleteBtn := makeButton("Delete Project", projectPane.removeActivateProject)
	clearBtn := makeButton("Clear Completed Tasks", clearCompletedTasks)

	deleteBtn.SetBackgroundColor(tcell.ColorRed)
	projectDetailPane = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(deleteBtn, 3, 1, false).
		AddItem(blankCell, 1, 1, false).
		AddItem(clearBtn, 3, 1, false).
		AddItem(blankCell, 0, 1, false)

	projectDetailPane.SetBorder(true).SetTitle("[::u]A[::-]ctions")
}

// @TODO - Move to tasks pane
func clearCompletedTasks() {
	count := 0
	for i, task := range taskPane.tasks {
		if task.Completed && taskRepo.Delete(&taskPane.tasks[i]) == nil {
			taskPane.list.RemoveItem(i)
			count++
		}
	}
	statusBar.showForSeconds(fmt.Sprintf("[yellow]%d tasks cleared!", count), 5)
}
