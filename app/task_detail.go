package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

var (
	taskName, taskDateDisplay *tview.TextView
	editorHint                *tview.TextView
	taskDate                  *tview.InputField
	taskDetailView            *femto.View
	taskStatusToggle          *tview.Button
	colorscheme               femto.Colorscheme
	blankCell                 = tview.NewTextView()
)

const dateLayoutISO = "2006-01-02"
const dateLayoutHuman = "02 Jan, Monday"

func prepareDetailPane() {
	taskName = tview.NewTextView().SetDynamicColors(true)

	prepareDetailsEditor()

	taskStatusToggle = makeButton("Complete", toggleActiveTaskStatus).SetLabelColor(tcell.ColorLightGray)

	toggleHint := tview.NewTextView().SetTextColor(tcell.ColorDimGray).SetText("<space> to toggle")

	editorLabel := tview.NewFlex().
		AddItem(tview.NewTextView().SetText("Task Not[::u]e[::-]:").SetDynamicColors(true), 0, 1, false).
		AddItem(makeButton("edit", func() { activateEditor() }), 6, 0, false)

	editorHint = tview.NewTextView().
		SetText(" e to edit, ↓↑ to scroll").
		SetTextColor(tcell.ColorDimGray)
	editorHelp := tview.NewFlex().
		AddItem(editorHint, 0, 1, false).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight).
			SetText("syntax:markdown theme:monakai").
			SetTextColor(tcell.ColorDimGray), 0, 1, false)

	taskDetailPane = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(taskName, 2, 1, true).
		AddItem(makeHorizontalLine(tcell.RuneS3, tcell.ColorGray), 1, 1, false).
		AddItem(blankCell, 1, 1, false).
		AddItem(makeDateRow(), 1, 1, true).
		AddItem(blankCell, 1, 1, false).
		AddItem(editorLabel, 1, 1, false).
		AddItem(taskDetailView, 15, 4, false).
		AddItem(editorHelp, 1, 1, false).
		AddItem(blankCell, 0, 1, false).
		AddItem(toggleHint, 1, 1, false).
		AddItem(taskStatusToggle, 3, 1, false)

	taskDetailPane.SetBorder(true).SetTitle("Task Detail")
}

func makeDateRow() *tview.Flex {
	taskDateDisplay = tview.NewTextView().SetDynamicColors(true)
	taskDate = makeLightTextInput("yyyy-mm-dd").
		SetLabel("Set:").
		SetLabelColor(tcell.ColorGray).
		SetFieldWidth(12).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				setTaskDate(parseDateInputOrCurrent(taskDate.GetText()).Unix(), true)
			case tcell.KeyEsc:
				setTaskDate(currentTask.DueDate, false)
			}
			app.SetFocus(taskDetailPane)
		})

	todaySelector := func() {
		setTaskDate(time.Now().Unix(), true)
	}

	nextDaySelector := func() {
		setTaskDate(parseDateInputOrCurrent(taskDate.GetText()).AddDate(0, 0, 1).Unix(), true)
	}

	prevDaySelector := func() {
		setTaskDate(parseDateInputOrCurrent(taskDate.GetText()).AddDate(0, 0, -1).Unix(), true)
	}

	return tview.NewFlex().
		AddItem(taskDateDisplay, 0, 2, true).
		AddItem(taskDate, 14, 0, true).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("today", todaySelector), 8, 1, false).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("+1", nextDaySelector), 4, 1, false).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("-1", prevDaySelector), 4, 1, false)
}

func setStatusToggle() {
	if currentTask.Completed {
		taskStatusToggle.SetLabel("Resume").SetBackgroundColor(tcell.ColorMaroon)
	} else {
		taskStatusToggle.SetLabel("Complete").SetBackgroundColor(tcell.ColorDarkGreen)
	}
}

func toggleActiveTaskStatus() {
	status := !currentTask.Completed
	if taskRepo.UpdateField(currentTask, "Completed", status) == nil {
		currentTask.Completed = status
		loadTask(currentTaskIdx)
		taskList.SetItemText(currentTaskIdx, makeTaskListingTitle(*currentTask), "")
	}
}

// Display Task date in detail pane, and update date if asked to
func setTaskDate(unixDate int64, update bool) {
	if update {
		currentTask.DueDate = unixDate
		if err := taskRepo.UpdateField(currentTask, "DueDate", unixDate); err != nil {
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
		taskDateDisplay.SetText(fmt.Sprintf("[::u]D[::-]ue: [%s]%s", color, humanDate))
		taskDate.SetText(due.Format(dateLayoutISO))
	} else {
		taskDate.SetText("")
		taskDateDisplay.SetText("[::u]D[::-]ue: [::d]Not Set")
	}
}

func prepareDetailsEditor() {
	taskDetailView = femto.NewView(makeBufferFromString(""))
	taskDetailView.SetRuntimeFiles(runtime.Files)

	// var colorscheme femto.Colorscheme
	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}
	taskDetailView.SetColorscheme(colorscheme)
	taskDetailView.SetBorder(true)
	taskDetailView.SetBorderColor(tcell.ColorLightSlateGray)

	taskDetailView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			currentTask.Details = taskDetailView.Buf.String()
			err := taskRepo.Update(currentTask)
			if err == nil {
				statusBar.showForSeconds("[lime]Saved task detail", 5)
			} else {
				statusBar.showForSeconds("[red]Could not save: "+err.Error(), 5)
			}

			deactivateEditor()
			return nil
		}

		return event
	})
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

func activateEditor() {
	taskDetailView.Readonly = false
	taskDetailView.SetBorderColor(tcell.ColorDarkOrange)
	editorHint.SetText(" Esc to save changes")
	app.SetFocus(taskDetailView)
}

func deactivateEditor() {
	taskDetailView.Readonly = true
	taskDetailView.SetBorderColor(tcell.ColorLightSlateGray)
	editorHint.SetText(" e to edit, ↓↑ to scroll")
	app.SetFocus(taskDetailPane)
}

func handleDetailPaneShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc:
		app.SetFocus(taskPane)
	case tcell.KeyDown:
		taskDetailView.ScrollDown(1)
	case tcell.KeyUp:
		taskDetailView.ScrollUp(1)
	case tcell.KeyRune:
		switch event.Rune() {
		case 'e':
			activateEditor()
		case 'd':
			app.SetFocus(taskDate)
		case ' ':
			toggleActiveTaskStatus()
		}
	}

	return event
}
