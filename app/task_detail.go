package main

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

var (
	taskName         *tview.TextView
	taskDate         *tview.InputField
	taskDetailView   *femto.View
	taskStatusToggle *tview.Button
	colorscheme      femto.Colorscheme
)

const dateLayoutISO = "2006-01-02"

func prepareDetailPane() {
	taskName = tview.NewTextView().SetDynamicColors(true)
	hr := makeHorizontalLine(tview.BoxDrawingsLightHorizontal)

	prepareDetailsEditor()

	taskStatusToggle = makeButton("Complete", func() {}).SetLabelColor(tcell.ColorLightGray)

	hint := tview.NewTextView().SetTextColor(tcell.ColorYellow).
		SetText("press Enter to save changes, Esc to ignore")

	detailPane = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(taskName, 2, 1, true).
		AddItem(hr, 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(makeDateRow(), 1, 1, true).
		AddItem(taskDetailView, 15, 4, false).
		AddItem(tview.NewTextView(), 1, 1, false).
		AddItem(hint, 1, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(taskStatusToggle, 3, 1, false)

	detailPane.SetBorder(true).SetTitle("Detail")

	// taskName is the default focus attracting child of detailPane
	taskName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyEsc:
			app.SetFocus(taskPane)
		case tcell.KeyDown:
			taskDetailView.ScrollDown(1)
		case tcell.KeyUp:
			taskDetailView.ScrollUp(1)
		case tcell.KeyRune:
			// switch event.Rune() {
			// case 'n':
			//     app.SetFocus(projectPane)
			// case 'e':
			//     if detailPane.HasFocus() {
			//         activateEditor()
			//     }
			// }
		}

		return event
	})
}

func makeDateRow() *tview.Flex {
	taskDate = makeLightTextInput("yyyy-mm-dd").
		SetLabel("Due Date: ").
		SetLabelColor(tcell.ColorGray).
		SetFieldWidth(12)

	todaySelector := func() {
		taskDate.SetText(time.Now().Format(dateLayoutISO))
	}

	nextDaySelector := func() {
		currentText := taskDate.GetText()
		if date, err := time.Parse(dateLayoutISO, currentText); err == nil {
			taskDate.SetText(date.AddDate(0, 0, 1).Format(dateLayoutISO))
		} else {
			taskDate.SetText(time.Now().AddDate(0, 0, 1).Format(dateLayoutISO))
		}
	}

	prevDaySelector := func() {
		currentText := taskDate.GetText()
		if date, err := time.Parse(dateLayoutISO, currentText); err == nil {
			taskDate.SetText(date.AddDate(0, 0, -1).Format(dateLayoutISO))
		} else {
			taskDate.SetText(time.Now().AddDate(0, 0, -1).Format(dateLayoutISO))
		}
	}

	return tview.NewFlex().
		AddItem(taskDate, 25, 2, true).
		AddItem(nil, 1, 0, false).
		AddItem(makeButton("today", todaySelector), 8, 1, false).
		AddItem(nil, 1, 0, false).
		AddItem(makeButton("+1", nextDaySelector), 4, 1, false).
		AddItem(nil, 1, 0, false).
		AddItem(makeButton("-1", prevDaySelector), 4, 1, false)
}

func setStatusToggle(idx int) {
	action := func(i int, label string, color tcell.Color, status bool) {
		taskStatusToggle.SetLabel(label).SetBackgroundColor(color)
		taskStatusToggle.SetSelectedFunc(func() {
			if taskRepo.UpdateField(currentTask, "Completed", status) == nil {
				currentTask.Completed = status
				loadTask(i)
				taskList.SetItemText(i, makeTaskListingTitle(*currentTask), "")
			}
		})
	}

	if currentTask.Completed {
		action(idx, "Resume", tcell.ColorMaroon, false)
	} else {
		action(idx, "Complete", tcell.ColorDarkGreen, true)
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
				showMessage("[lime]Saved task detail")
			} else {
				showMessage("[red]Could not save: " + err.Error())
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
	app.SetFocus(taskDetailView)
}

func deactivateEditor() {
	taskDetailView.Readonly = true
	taskDetailView.SetBorderColor(tcell.ColorLightSlateGray)
	app.SetFocus(detailPane)
}
