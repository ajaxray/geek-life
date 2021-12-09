package main

import (
	"fmt"
	"sort"
	"time"
	"unicode"

	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
)

// TaskPane displays tasks of current TaskList or Project
type TaskPane struct {
	*tview.Flex
	list       *tview.List
	tasks      []model.Task
	activeTask *model.Task

	newTask     *tview.InputField
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository
	hint        *tview.TextView
}

// NewTaskPane initializes and configures a TaskPane
func NewTaskPane(projectRepo repository.ProjectRepository, taskRepo repository.TaskRepository) *TaskPane {
	pane := TaskPane{
		Flex:        tview.NewFlex().SetDirection(tview.FlexRow),
		list:        tview.NewList().ShowSecondaryText(false),
		newTask:     makeLightTextInput("+[New Task]"),
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
		hint:        tview.NewTextView().SetTextColor(tcell.ColorYellow).SetTextAlign(tview.AlignCenter),
	}

	pane.list.SetSelectedBackgroundColor(tcell.ColorDarkBlue)
	pane.list.SetDoneFunc(func() {
		app.SetFocus(projectPane)
	})

	pane.newTask.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			name := pane.newTask.GetText()
			if len(name) < 3 {
				statusBar.showForSeconds("[red::]Task title should be at least 3 character", 5)
				return
			}

			task, err := taskRepo.Create(*projectPane.GetActiveProject(), name, "", "", 0)
			if err != nil {
				statusBar.showForSeconds("[red::]Could not create Task:"+err.Error(), 5)
				return
			}

			pane.tasks = append(pane.tasks, task)
			pane.addTaskToList(len(pane.tasks) - 1)
			pane.newTask.SetText("")
			statusBar.showForSeconds("[yellow::]Task created. Add another task or press Esc.", 5)
		case tcell.KeyEsc:
			app.SetFocus(pane)
		}
	})

	pane.
		AddItem(pane.list, 0, 1, true).
		AddItem(pane.hint, 0, 1, false)

	pane.SetBorder(true).SetTitle("[::u]T[::-]asks")
	pane.setHintMessage()

	return &pane
}

// ClearList removes all items from TaskPane
func (pane *TaskPane) ClearList() {
	pane.list.Clear()
	pane.tasks = nil
	pane.activeTask = nil

	pane.RemoveItem(pane.newTask)
}

// SetList Sets a list of tasks to be displayed
func (pane *TaskPane) SetList(tasks []model.Task) {
	pane.ClearList()
	pane.tasks = tasks

	for i := range pane.tasks {
		pane.addTaskToList(i)
	}
}

func (pane *TaskPane) addTaskToList(i int) *tview.List {
	return pane.list.AddItem(makeTaskListingTitle(pane.tasks[i]), "", 0, func(taskidx int) func() {
		return func() { taskPane.ActivateTask(taskidx) }
	}(i))
}

func (pane *TaskPane) handleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch unicode.ToLower(event.Rune()) {
	case 'j':
		pane.list.SetCurrentItem(pane.list.GetCurrentItem() + 1)
		return nil
	case 'k':
		pane.list.SetCurrentItem(pane.list.GetCurrentItem() - 1)
		return nil
	case 'h':
		app.SetFocus(projectPane)
		return nil
	case 'n':
		app.SetFocus(pane.newTask)
		return nil
	}

	return event
}

// LoadProjectTasks loads tasks of a project in taskPane
func (pane *TaskPane) LoadProjectTasks(project model.Project) {
	var tasks []model.Task
	var err error

	if tasks, err = taskRepo.GetAllByProject(project); err != nil && err != storm.ErrNotFound {
		statusBar.showForSeconds("[red::]Error: "+err.Error(), 5)
	} else {
		pane.SetList(tasks)
	}

	pane.RemoveItem(pane.hint)
	pane.AddItem(pane.newTask, 1, 0, false)
}

// LoadDynamicList loads tasks based on logic key
func (pane *TaskPane) LoadDynamicList(logic string) {
	var tasks []model.Task
	var err error

	today := toDate(time.Now())
	zeroTime := time.Time{}
	rangeDesc := ""

	switch logic {
	case "today":
		tasks, err = pane.taskRepo.GetAllByDateRange(zeroTime, today)
		rangeDesc = "Today (and overdue)"

	case "tomorrow":
		tomorrow := today.AddDate(0, 0, 1)
		tasks, err = pane.taskRepo.GetAllByDate(tomorrow)
		rangeDesc = "Tomorrow"

	case "upcoming":
		week := today.Add(7 * 24 * time.Hour)
		tasks, err = pane.taskRepo.GetAllByDateRange(today, week)
		rangeDesc = "Upcoming (next 7 days)"

	case "unscheduled":
		tasks, err = pane.taskRepo.GetAllByDate(zeroTime)
		rangeDesc = "Unscheduled (task with no due date) "
	}

	projectPane.activeProject = nil
	taskPane.ClearList()

	if err == storm.ErrNotFound {
		statusBar.showForSeconds("[yellow]No Task in list - "+rangeDesc, 5)
		pane.SetList(tasks)
	} else if err != nil {
		statusBar.showForSeconds("[red]Error: "+err.Error(), 5)
	} else {
		sort.Slice(tasks, func(i, j int) bool { return tasks[i].ProjectID < tasks[j].ProjectID })
		pane.SetList(tasks)
		app.SetFocus(taskPane)
		contents.RemoveItem(taskDetailPane)

		statusBar.showForSeconds("[yellow] Displaying tasks of "+rangeDesc, 5)
	}

	pane.RemoveItem(pane.hint)
	removeThirdCol()
}

// ActivateTask marks a task as currently active and loads in TaskDetailPane
func (pane *TaskPane) ActivateTask(idx int) {
	removeThirdCol()
	pane.activeTask = &pane.tasks[idx]
	taskDetailPane.SetTask(pane.activeTask)

	contents.AddItem(taskDetailPane, 0, 3, false)

}

// ClearCompletedTasks removes tasks from current list that are in completed state
func (pane *TaskPane) ClearCompletedTasks() {
	count := 0
	for i, task := range pane.tasks {
		if task.Completed && pane.taskRepo.Delete(&pane.tasks[i]) == nil {
			pane.list.RemoveItem(i)
			count++
		}
	}

	statusBar.showForSeconds(fmt.Sprintf("[yellow]%d tasks cleared!", count), 5)
}

// ReloadCurrentTask Loads the current task - in Task details and listing
func (pane *TaskPane) ReloadCurrentTask() {
	pane.list.SetItemText(pane.list.GetCurrentItem(), makeTaskListingTitle(*pane.activeTask), "")
	taskDetailPane.SetTask(pane.activeTask)
}

func (pane TaskPane) setHintMessage() {
	if len(projectPane.projects) == 0 {
		pane.hint.SetText("Welcome to the organized life!\n------------------------------\n Create TaskList/Project at the bottom of Projects pane.\n (Press p,n) \n\nHelp - https://bit.ly/cli-task")
	} else {
		pane.hint.SetText("Select a TaskList/Project (Press Enter) to load tasks.\nOr create a new Project (Press p,n).\n\nHelp - https://bit.ly/cli-task")
	}

	// Add: For help - https://bit.ly/cli-task
}
