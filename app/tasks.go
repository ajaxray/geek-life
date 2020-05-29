package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
)

var (
	tasks       []model.Task
	currentTask *model.Task
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
				task, err := taskRepo.Create(currentProject, newTask.GetText(), "", "", time.Now().Unix())
				if err != nil {
					showMessage("[red::]Could not create Task:" + err.Error())
				}

				taskList.AddItem(task.Title, "", 0, nil)
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

func loadTask(idx int) {
	contents.RemoveItem(detailPane)
	currentTask = &tasks[idx]

	taskName.SetText(fmt.Sprintf("[%s::b]# %s", getTaskTitleColor(*currentTask), currentTask.Title))
	taskDetailView.Buf = makeBufferFromString(currentTask.Details)
	taskDetailView.SetColorscheme(colorscheme)
	taskDetailView.Start()

	taskDate.SetText("")
	if currentTask.DueDate != 0 {
		taskDate.SetText(time.Unix(currentTask.DueDate, 0).Format(dateLayoutISO))
	}

	contents.AddItem(detailPane, 0, 3, false)
	setStatusToggle(idx)
	deactivateEditor()
}

func getTaskTitleColor(task model.Task) string {
	colorName := "whitesmoke"
	if task.Completed {
		colorName = "lime"
	} else if task.DueDate != 0 && task.DueDate < time.Now().Unix() {
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
