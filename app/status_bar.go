package main

import (
	"time"

	"github.com/rivo/tview"
)

type StatusBar struct {
	*tview.Pages
	message   *tview.TextView
	container *tview.Application
}

// Name of page keys
const (
	defaultPage = "default"
	messagePage = "message"
)

func makeStatusBar(app *tview.Application) *StatusBar {
	statusBar := StatusBar{
		Pages:     tview.NewPages(),
		message:   tview.NewTextView().SetDynamicColors(true).SetText("Loading..."),
		container: app,
	}

	statusBar.AddPage(messagePage, statusBar.message, true, true)
	statusBar.AddPage(defaultPage,
		tview.NewGrid(). // Content will not be modified, So, no need to declare explicitly
					SetColumns(0, 0, 0, 0).
					SetRows(0).
					AddItem(tview.NewTextView().SetText("Navigate List: ↓/↑"), 0, 0, 1, 1, 0, 0, false).
					AddItem(tview.NewTextView().SetText("New Task/Project: n").SetTextAlign(tview.AlignCenter), 0, 1, 1, 1, 0, 0, false).
					AddItem(tview.NewTextView().SetText("Step back: Esc").SetTextAlign(tview.AlignCenter), 0, 2, 1, 1, 0, 0, false).
					AddItem(tview.NewTextView().SetText("Quit: Ctrl+C").SetTextAlign(tview.AlignRight), 0, 3, 1, 1, 0, 0, false),
		true,
		true,
	)

	return &statusBar
}

func (bar *StatusBar) showForSeconds(message string, timeout int) {
	if bar.container == nil {
		return
	}

	bar.message.SetText(message)
	bar.SwitchToPage(messagePage)

	go func() {
		time.Sleep(time.Second * time.Duration(timeout))
		bar.container.QueueUpdateDraw(func() {
			bar.SwitchToPage(defaultPage)
		})
	}()
}
