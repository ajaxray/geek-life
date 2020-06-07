package main

import (
	"fmt"
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func prepareProjectPane() {
	projectList = tview.NewList().ShowSecondaryText(false)
	loadProjectList()

	newProject = makeLightTextInput("+[New Project]").
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				if len(newProject.GetText()) < 3 {
					statusBar.showForSeconds("[red::]Project name should be at least 3 character", 5)
					return
				}

				project, err := projectRepo.Create(newProject.GetText(), "")
				if err != nil {
					statusBar.showForSeconds("[red::]Failed to create Project:"+err.Error(), 5)
				} else {
					statusBar.showForSeconds(fmt.Sprintf("[yellow::]Project %s created. Press n to start adding new tasks.", newProject.GetText()), 5)
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

func loadProjectList() {
	var err error
	projects, err = projectRepo.GetAll()
	if err != nil {
		statusBar.showForSeconds("Could not load Projects: "+err.Error(), 5)
		return
	}

	projectList.AddItem("[::d]Dynamic Lists", "", 0, nil)
	projectList.AddItem("[::d]"+strings.Repeat(string(tcell.RuneS3), 25), "", 0, nil)
	projectList.AddItem("- Today", "", 0, yetToImplement("Today's Tasks"))
	projectList.AddItem("- Upcoming", "", 0, yetToImplement("Upcoming Tasks"))
	projectList.AddItem("- No Due Date", "", 0, yetToImplement("Unscheduled Tasks"))
	projectList.AddItem("", "", 0, nil)
	projectList.AddItem("[::d]Projects", "", 0, nil)
	projectList.AddItem("[::d]"+strings.Repeat(string(tcell.RuneS3), 25), "", 0, nil)

	for i := range projects {
		addProjectToList(i, false)
	}

	projectList.SetCurrentItem(6) // Select Projects, as dynamic lists are not ready
}

func addProjectToList(i int, selectItem bool) {
	// To avoid overriding of loop variables - https://www.calhoun.io/gotchas-and-common-mistakes-with-closures-in-go/
	projectList.AddItem("- "+projects[i].Title, "", 0, func(idx int) func() {
		return func() { loadProject(idx) }
	}(i))

	if selectItem {
		projectList.SetCurrentItem(projectList.GetItemCount() - 1)
		loadProject(i)
	}
}

func loadProject(idx int) {
	currentProject = &projects[idx]
	taskList.Clear()
	app.SetFocus(taskPane)
	var err error

	if tasks, err = taskRepo.GetAllByProject(*currentProject); err != nil && err != storm.ErrNotFound {
		statusBar.showForSeconds("[red::]Error: "+err.Error(), 5)
	}

	for i, task := range tasks {
		addTaskToList(task, i)
	}

	removeThirdCol()
	projectDetailPane.SetTitle("[::b]" + currentProject.Title)
	contents.AddItem(projectDetailPane, 25, 0, false)
}

func handleProjectPaneShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'n':
		app.SetFocus(newProject)
	}

	return event
}
