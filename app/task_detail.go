package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
	"github.com/ajaxray/geek-life/util"
)

const dateLayoutISO = "2006-01-02"
const dateLayoutHuman = "02 Jan, Monday"

// TaskDetailPane displays detailed info of a Task
type TaskDetailPane struct {
	*tview.Flex
	header           *TaskDetailHeader
	taskDateDisplay  *tview.TextView
	editorHint       *tview.TextView
	taskDate         *tview.InputField
	taskStatusToggle *tview.Button
	taskDetailView   *femto.View
	colorScheme      femto.Colorscheme
	taskRepo         repository.TaskRepository
	task             *model.Task
}

// NewTaskDetailPane initializes and configures a TaskDetailPane
func NewTaskDetailPane(taskRepo repository.TaskRepository) *TaskDetailPane {
	pane := TaskDetailPane{
		Flex:             tview.NewFlex().SetDirection(tview.FlexRow),
		header:           NewTaskDetailHeader(taskRepo),
		taskDateDisplay:  tview.NewTextView().SetDynamicColors(true),
		taskStatusToggle: makeButton("Complete", nil).SetLabelColor(tcell.ColorLightGray),
		taskRepo:         taskRepo,
	}

	pane.prepareDetailsEditor()

	toggleHint := tview.NewTextView().SetTextColor(tcell.ColorDimGray).SetText("<space> to toggle")
	pane.taskStatusToggle.SetSelectedFunc(pane.toggleTaskStatus)

	pane.editorHint = tview.NewTextView().SetText(" e to edit, ↓↑ to scroll").SetTextColor(tcell.ColorDimGray)

	// Prepare static (no external interaction) elements
	editorLabel := tview.NewFlex().
		AddItem(tview.NewTextView().SetText("Task Not[::u]e[::-]:").SetDynamicColors(true), 0, 1, false).
		AddItem(makeButton("edit", func() { pane.activateEditor() }), 6, 0, false)
	editorHelp := tview.NewFlex().
		AddItem(pane.editorHint, 0, 1, false).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight).
			SetText("syntax:markdown theme:monakai").
			SetTextColor(tcell.ColorDimGray), 0, 1, false)

	pane.
		AddItem(pane.header, 4, 1, true).
		AddItem(blankCell, 1, 1, false).
		AddItem(pane.makeDateRow(), 1, 1, true).
		AddItem(blankCell, 1, 1, false).
		AddItem(editorLabel, 1, 1, false).
		AddItem(pane.taskDetailView, 15, 4, false).
		AddItem(editorHelp, 1, 1, false).
		AddItem(blankCell, 0, 1, false).
		AddItem(toggleHint, 1, 1, false).
		AddItem(pane.taskStatusToggle, 3, 1, false)

	pane.SetBorder(true).SetTitle("Task Detail")

	return &pane
}

func (td *TaskDetailPane) makeDateRow() *tview.Flex {

	td.taskDate = makeLightTextInput("yyyy-mm-dd").
		SetLabel("Set:").
		SetLabelColor(tcell.ColorGray).
		SetFieldWidth(12).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				date := parseDateInputOrCurrent(td.taskDate.GetText())
				td.setTaskDate(date.Unix(), true)
			case tcell.KeyEsc:
				td.setTaskDate(td.task.DueDate, false)
			}
			app.SetFocus(td)
		})

	todaySelector := func() {
		td.setTaskDate(parseDateInputOrCurrent("").Unix(), true)
	}

	nextDaySelector := func() {
		td.setTaskDate(parseDateInputOrCurrent(td.taskDate.GetText()).AddDate(0, 0, 1).Unix(), true)
	}

	prevDaySelector := func() {
		td.setTaskDate(parseDateInputOrCurrent(td.taskDate.GetText()).AddDate(0, 0, -1).Unix(), true)
	}

	return tview.NewFlex().
		AddItem(td.taskDateDisplay, 0, 2, true).
		AddItem(td.taskDate, 14, 0, true).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("today", todaySelector), 8, 1, false).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("+1", nextDaySelector), 4, 1, false).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("-1", prevDaySelector), 4, 1, false)
}

func (td *TaskDetailPane) updateToggleDisplay() {
	if td.task.Completed {
		td.taskStatusToggle.SetLabel("Resume").SetBackgroundColor(tcell.ColorMaroon)
	} else {
		td.taskStatusToggle.SetLabel("Complete").SetBackgroundColor(tcell.ColorDarkGreen)
	}
}

func (td *TaskDetailPane) toggleTaskStatus() {
	status := !td.task.Completed
	if taskRepo.UpdateField(td.task, "Completed", status) == nil {
		td.task.Completed = status
		taskPane.ReloadCurrentTask()
	}
}

// Display Task date in detail pane, and update date if asked to
func (td *TaskDetailPane) setTaskDate(unixDate int64, update bool) {
	if update {
		td.task.DueDate = unixDate
		if err := td.taskRepo.UpdateField(td.task, "DueDate", unixDate); err != nil {
			statusBar.showForSeconds("Could not update due date: "+err.Error(), 5)
			return
		}
	}

	if unixDate != 0 {
		due := time.Unix(unixDate, 0)
		color := "white"
		humanDate := due.Format(dateLayoutHuman)

		if due.Before(time.Now()) {
			color = "red"
		}
		td.taskDateDisplay.SetText(fmt.Sprintf("[::u]D[::-]ue: [%s]%s", color, humanDate))
		td.taskDate.SetText(due.Format(dateLayoutISO))
	} else {
		td.taskDate.SetText("")
		td.taskDateDisplay.SetText("[::u]D[::-]ue: [::d]Not Set")
	}
}

func (td *TaskDetailPane) prepareDetailsEditor() {

	td.taskDetailView = femto.NewView(makeBufferFromString(""))
	td.taskDetailView.SetRuntimeFiles(runtime.Files)

	// var colorScheme femto.Colorscheme
	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			td.colorScheme = femto.ParseColorscheme(string(data))
		}
	}

	td.taskDetailView.SetColorscheme(td.colorScheme)
	td.taskDetailView.SetBorder(true)
	td.taskDetailView.SetBorderColor(tcell.ColorLightSlateGray)

	td.taskDetailView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			td.updateTaskNote(td.taskDetailView.Buf.String())
			td.deactivateEditor()
			return nil
		}

		return event
	})
}

func (td *TaskDetailPane) updateTaskNote(note string) {
	td.task.Details = note
	err := taskRepo.Update(td.task)
	if err == nil {
		statusBar.showForSeconds("[lime]Saved task detail", 5)
	} else {
		statusBar.showForSeconds("[red]Could not save: "+err.Error(), 5)
	}
}

func makeBufferFromString(content string) *femto.Buffer {
	buff := femto.NewBufferFromString(content, "")
	// taskDetail.Settings["ruler"] = false
	buff.Settings["filetype"] = "markdown"
	buff.Settings["keepautoindent"] = true
	buff.Settings["statusline"] = false
	buff.Settings["softwrap"] = true
	buff.Settings["scrollbar"] = true

	return buff
}

func (td *TaskDetailPane) activateEditor() {
	td.taskDetailView.Readonly = false
	td.taskDetailView.SetBorderColor(tcell.ColorDarkOrange)
	td.editorHint.SetText(" Esc to save changes")
	app.SetFocus(td.taskDetailView)
}

func (td *TaskDetailPane) deactivateEditor() {
	td.taskDetailView.Readonly = true
	td.taskDetailView.SetBorderColor(tcell.ColorLightSlateGray)
	td.editorHint.SetText(" e to edit, ↓↑ to scroll")
	app.SetFocus(td)
}

func (td *TaskDetailPane) editInExternalEditor() {

	tmpFileName, err := writeToTmpFile(td.task.Details)
	if err != nil {
		statusBar.showForSeconds("[red::]Failed to create tmp file. Try in-app editing by pressing i", 5)
		return
	}

	var messageToShow, updatedContent string
	app.Suspend(func() {
		cmd := exec.Command(util.GetEnvStr("EDITOR", "vim"), tmpFileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			messageToShow = "[red::]Failed to save content. Try in-app editing by pressing e"
			return
		}

		if content, readErr := ioutil.ReadFile(tmpFileName); readErr == nil {
			updatedContent = string(content)
		} else {
			messageToShow = "[red::]Failed to load external editing. Try in-app editing by pressing e"
		}
	})

	if messageToShow != "" {
		statusBar.showForSeconds(messageToShow, 10)
	}

	if updatedContent != "" {
		td.updateTaskNote(updatedContent)
		td.SetTask(td.task)
	}

	// @TODO: Mouse events not being captured after returning from Suspend - fix it
	app.EnableMouse(true)

	_ = os.Remove(tmpFileName)

	// app.SetFocus(td)
}

// writeToTmpFile writes given content to a tmpFile and returns the filename
func writeToTmpFile(content string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "geek_life_task_note_*.md")
	if err != nil {
		return "", err
	}
	fileName := tmpFile.Name()

	if err = ioutil.WriteFile(fileName, []byte(content), 0777); err != nil {
		return "", err
	}

	return fileName, tmpFile.Close()
}

func (td *TaskDetailPane) handleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc:
		app.SetFocus(taskPane)
	case tcell.KeyDown:
		td.taskDetailView.ScrollDown(1)
	case tcell.KeyUp:
		td.taskDetailView.ScrollUp(1)
	case tcell.KeyRune:
		switch event.Rune() {
		case 'e':
			td.activateEditor()
		case 'v':
			td.editInExternalEditor()
		case 'd':
			app.SetFocus(td.taskDate)
		case 'r':
			td.header.ShowRename()
		case ' ':
			td.toggleTaskStatus()
		}
	}

	return event
}

// SetTask sets a Task to be displayed
func (td *TaskDetailPane) SetTask(task *model.Task) {
	td.task = task

	td.header.SetTask(task)
	td.taskDetailView.Buf = makeBufferFromString(td.task.Details)
	td.taskDetailView.SetColorscheme(td.colorScheme)
	td.taskDetailView.Start()
	td.setTaskDate(td.task.DueDate, false)
	td.updateToggleDisplay()
	td.deactivateEditor()
}
