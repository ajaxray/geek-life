package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
)

var (
	tasks          []model.Task
	currentTask    *model.Task
	currentTaskIdx int
)

func prepareTaskPane() {
	taskList = tview.NewList().ShowSecondaryText(false)
	taskList.SetDoneFunc(func() {
		app.SetFocus(projectPane)
	})

	newTask = makeLightTextInput("+[New Task]").
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				if len(newTask.GetText()) < 3 {
					showMessage("[red::]Task title should be at least 3 character")
					return
				}

				task, err := taskRepo.Create(*currentProject, newTask.GetText(), "", "", 0)
				if err != nil {
					showMessage("[red::]Could not create Task:" + err.Error())
				}

				tasks = append(tasks, task)
				addTaskToList(task, len(tasks)-1)
				newTask.SetText("")
			case tcell.KeyEsc:
				app.SetFocus(taskPane)
			}

		})

	taskPane = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(taskList, 0, 1, true).
		AddItem(newTask, 1, 0, false)

	taskPane.SetBorder(true).SetTitle("[::u]T[::-]asks")
}

func addTaskToList(task model.Task, i int) *tview.List {
	return taskList.AddItem(makeTaskListingTitle(task), "", 0, func(taskidx int) func() {
		return func() { loadTask(taskidx) }
	}(i))
}

func loadTask(idx int) {
	removeThirdCol()
	currentTaskIdx = idx
	currentTask = &tasks[currentTaskIdx]

	taskName.SetText(fmt.Sprintf("[%s::b]# %s", getTaskTitleColor(*currentTask), currentTask.Title))
	taskDetailView.Buf = makeBufferFromString(currentTask.Details)
	taskDetailView.SetColorscheme(colorscheme)
	taskDetailView.Start()
	setTaskDate(currentTask.DueDate, false)
	setStatusToggle()

	contents.AddItem(taskDetailPane, 0, 3, false)
	deactivateEditor()
}

func removeThirdCol() {
	contents.RemoveItem(taskDetailPane)
	contents.RemoveItem(projectDetailPane)
}

func getTaskTitleColor(task model.Task) string {
	colorName := "olive"
	if task.Completed {
		colorName = "lime"
	} else if task.DueDate != 0 && task.DueDate < time.Now().Truncate(24*time.Hour).Unix() {
		colorName = "red"
	}

	return colorName
}

func makeTaskListingTitle(task model.Task) string {
	checkbox := "[ []"
	if task.Completed {
		checkbox = "[x[]"
	}
	return fmt.Sprintf("[%s]%s %s", getTaskTitleColor(task), checkbox, task.Title)
}

func handleTaskPaneShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'n':
		app.SetFocus(newTask)
	}

	return event
}
