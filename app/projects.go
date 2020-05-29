package main

import (
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func prepareProjectPane() {
	var err error
	projects, err = projectRepo.GetAll()
	if err != nil {
		showMessage("Could not load Projects: " + err.Error())
	}

	projectList = tview.NewList().ShowSecondaryText(false)

	for i := range projects {
		addProjectToList(i, false)
	}

	newProject = makeLightTextInput("+[New Project]").
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				project, err := projectRepo.Create(newProject.GetText(), "")
				if err != nil {
					showMessage("[red::]Failed to create Project:" + err.Error())
				} else {
					showMessage(fmt.Sprintf("[green::]Project %s created. Press n to start adding new tasks.", newProject.GetText()))
					projects = append(projects, project)
					addProjectToList(len(projects)-1, true)
					newProject.SetText("")
				}
			case tcell.KeyEsc:
				app.SetFocus(projectPane)
			}
		})

	projectPane = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(projectList, 0, 1, true).
		AddItem(newProject, 1, 0, false)

	projectPane.SetBorder(true).SetTitle("[::u]P[::-]rojects")
}

func addProjectToList(i int, selectItem bool) {
	// To avoid overriding of loop variables - https://www.calhoun.io/gotchas-and-common-mistakes-with-closures-in-go/
	projectList.AddItem("- "+projects[i].Title, "", 0, func(idx int) func() {
		return func() { loadProject(idx) }
	}(i))

	if selectItem {
		projectList.SetCurrentItem(i)
		loadProject(i)
	}
}

func loadProject(idx int) {
	currentProject = projects[idx]
	taskList.Clear()
	app.SetFocus(taskPane)
	var err error

	if tasks, err = taskRepo.GetAllByProject(currentProject); err != nil && err != storm.ErrNotFound {
		showMessage("[red::]Error: " + err.Error())
	}

	for i, task := range tasks {
		taskList.AddItem(makeTaskListingTitle(task), "", 0, func(taskidx int) func() {
			return func() { loadTask(taskidx) }
		}(i))
	}

	contents.RemoveItem(detailPane)
}
