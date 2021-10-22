package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
)

// TaskDetailHeader displays Task title and relevant action in TaskDetail pane
type TaskDetailHeader struct {
	*tview.Flex
	pages      *tview.Pages
	taskName   *tview.TextView
	taskRepo   repository.TaskRepository
	task       *model.Task
	renameText *tview.InputField
}

// NewTaskDetailHeader initializes and configures a TaskDetailHeader
func NewTaskDetailHeader(taskRepo repository.TaskRepository) *TaskDetailHeader {
	header := TaskDetailHeader{
		Flex:       tview.NewFlex().SetDirection(tview.FlexRow),
		pages:      tview.NewPages(),
		taskName:   tview.NewTextView().SetDynamicColors(true),
		taskRepo:   taskRepo,
		renameText: makeLightTextInput("Task title"),
	}

	header.pages.AddPage("title", header.taskName, true, true)
	header.pages.AddPage("rename", header.renameText, true, false)
	header.bindRenameEvent()

	buttons := tview.NewFlex().
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorDimGray).SetText("r = Rename"), 0, 1, false).
		AddItem(blankCell, 0, 1, false).
		AddItem(makeButton("[::ub]r[::-]ename", func() { header.ShowRename() }), 8, 0, false).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("e[::ub]x[::-]port", func() { taskDetailPane.Export() }), 8, 0, false)

	header.
		AddItem(header.pages, 1, 1, true).
		AddItem(blankCell, 1, 0, true).
		AddItem(buttons, 1, 1, false).
		AddItem(makeHorizontalLine(tcell.RuneS3, tcell.ColorGray), 1, 1, false)

	return &header
}

func (header *TaskDetailHeader) bindRenameEvent() *tview.InputField {
	return header.renameText.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			name := header.renameText.GetText()
			if len(name) < 3 {
				statusBar.showForSeconds("[red::]Task title should be at least 3 character", 5)
				return
			}

			if err := header.taskRepo.UpdateField(header.task, "Title", name); err != nil {
				statusBar.showForSeconds("Could not update Task Title: "+err.Error(), 5)
			} else {
				header.task.Title = name
				_ = header.taskRepo.UpdateField(header.task, "ModifiedAt", time.Now().Unix())
				statusBar.showForSeconds("[yellow::]Task Title Updated.", 5)
			}

			header.pages.SwitchToPage("title")
			taskPane.ReloadCurrentTask()
		case tcell.KeyEsc:
			header.pages.SwitchToPage("title")
		}
	})
}

// SetTask set a task to the header
func (header *TaskDetailHeader) SetTask(task *model.Task) {
	header.task = task
	header.taskName.SetText(task.Title)
	header.taskName.SetText(fmt.Sprintf("[%s::b]# %s", getTaskTitleColor(*task), task.Title))
}

// ShowRename activate edit option of task title
func (header *TaskDetailHeader) ShowRename() {
	header.renameText.SetText(header.task.Title)
	header.pages.SwitchToPage("rename")
	app.SetFocus(header.renameText)
}
