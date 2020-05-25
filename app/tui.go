package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
	repo "github.com/ajaxray/geek-life/repository/storm"
	"github.com/ajaxray/geek-life/util"
)

var (
	app                        *tview.Application
	newProject, newTask        *tview.InputField
	projectList, taskList      *tview.List
	projectPane, taskPane      *tview.Flex
	statusBar                  *tview.Pages
	message                    *tview.TextView
	shortcutsPage, messagePage string = "shortcuts", "message"

	db          *storm.DB
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository

	projects       []model.Project
	currentProject model.Project

	tasks       []model.Task
	currentTask model.Task
)

func main() {
	app = tview.NewApplication()

	db = util.ConnectStorm()
	defer db.Close()

	projectRepo = repo.NewProjectRepository(db)
	taskRepo = repo.NewTaskRepository(db)

	titleText := tview.NewTextView().SetText("[lime::b]Geek-life [::-]- life management for geeks!").SetDynamicColors(true)
	cloudStatus := tview.NewTextView().SetText("[::d]Cloud Sync: off").SetTextAlign(tview.AlignRight).SetDynamicColors(true)

	prepareStatusBar()

	titleBar := tview.NewFlex().
		AddItem(titleText, 0, 2, false).
		AddItem(cloudStatus, 0, 1, false)

	projectPane = makeProjectPane()
	taskPane = makeTaskPane()

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBar, 2, 1, false).
		AddItem(tview.NewFlex().
			AddItem(projectPane, 25, 1, true).
			AddItem(taskPane, 0, 2, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Details"), 0, 3, false),
			0, 2, true).
		AddItem(statusBar, 1, 1, false)

	setKeyboardShortcuts(projectPane, taskPane)

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func makeProjectPane() *tview.Flex {
	var err error
	projects, err = projectRepo.GetAll()
	util.FatalIfError(err, "Could not load Projects")

	projectList = tview.NewList()

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

	projectBar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(projectList, 0, 1, true).
		AddItem(newProject, 1, 0, false)

	projectBar.SetBorder(true).SetTitle("Projects (p)")

	return projectBar
}

func addProjectToList(i int, selectItem bool) {
	// To avoid overriding of loop variables - https://www.calhoun.io/gotchas-and-common-mistakes-with-closures-in-go/
	projectList.AddItem(projects[i].Title, "", 0, func(idx int) func() {
		return func() { loadProject(idx) }
	}(i))

	if selectItem {
		projectList.SetCurrentItem(i)
		loadProject(i)
	}
}

func makeTaskPane() *tview.Flex {
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

	taskPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(taskList, 0, 1, true).
		AddItem(newTask, 1, 0, false)

	taskPane.SetBorder(true).SetTitle("Tasks (t)")

	return taskPane
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
		taskList.AddItem(task.Title, "", 0, func(taskidx int) func() {
			return func() { loadTask(taskidx) }
		}(i))
	}
}

func loadTask(idx int) {
	currentTask = tasks[idx]
	// taskList.Clear()
	// app.SetFocus(taskPane)
}

func makeLightTextInput(placeholder string) *tview.InputField {
	return tview.NewInputField().
		SetPlaceholder(placeholder).
		SetPlaceholderTextColor(tcell.ColorLightSlateGray).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorGray)
}

func ignoreKeyEvt() bool {
	return reflect.TypeOf(app.GetFocus()).String() == "*tview.InputField"
}

func setKeyboardShortcuts(projectPane *tview.Flex, taskPane *tview.Flex) *tview.Application {
	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvt() {
			return event
		}
		switch event.Rune() {
		case 'p':
			app.SetFocus(projectPane)
		case 't':
			app.SetFocus(taskPane)
		case 'n':
			if projectPane.HasFocus() {
				app.SetFocus(newProject)
			} else if taskPane.HasFocus() {
				app.SetFocus(newTask)
			}
		}

		return event
	})
}

func showMessage(text string) {
	message.SetText(text)
	statusBar.SwitchToPage(messagePage)

	go func() {
		app.QueueUpdateDraw(func() {
			time.Sleep(time.Second * 5)
			statusBar.SwitchToPage(shortcutsPage)
		})
	}()
}

func prepareStatusBar() {
	statusBar = tview.NewPages()

	message = tview.NewTextView().SetDynamicColors(true).SetText("Loading...")
	statusBar.AddPage(messagePage, message, true, true)

	statusBar.AddPage(shortcutsPage,
		tview.NewGrid().
			SetColumns(0, 0, 0, 0).
			SetRows(0).
			AddItem(tview.NewTextView().SetText("Shortcuts: Alt+.(dot)"), 0, 0, 1, 1, 0, 0, false).
			AddItem(tview.NewTextView().SetText("New Project: n").SetTextAlign(tview.AlignCenter), 0, 1, 1, 1, 0, 0, false).
			AddItem(tview.NewTextView().SetText("New Task: t").SetTextAlign(tview.AlignCenter), 0, 2, 1, 1, 0, 0, false).
			AddItem(tview.NewTextView().SetText("Quit: Ctrl+C").SetTextAlign(tview.AlignRight), 0, 3, 1, 1, 0, 0, false),
		true,
		true,
	)
}
