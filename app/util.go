package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/util"
)

func makeHorizontalLine(lineChar rune, color tcell.Color) *tview.TextView {
	hr := tview.NewTextView()
	hr.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		// Draw a horizontal line across the middle of the box.
		style := tcell.StyleDefault.Foreground(color).Background(tcell.ColorBlack)
		centerY := y + height/2
		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, centerY, lineChar, nil, style)
		}

		// Space for other content.
		return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
	})

	return hr
}

func makeLightTextInput(placeholder string) *tview.InputField {
	return tview.NewInputField().
		SetPlaceholder(placeholder).
		SetPlaceholderTextColor(tcell.ColorYellow).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorGray)
}

// If input text is a valid date, parse it. Or get current date
func parseDateInputOrCurrent(inputText string) time.Time {
	if date, err := time.Parse(dateLayoutISO, inputText); err == nil {
		return date
	} else {
		return time.Now()
	}
}

func showMessage(text string) {
	message.SetText(text)
	statusBar.SwitchToPage(messagePage)

	go func() {
		time.Sleep(time.Second * 5)
		app.QueueUpdateDraw(func() {
			statusBar.SwitchToPage(shortcutsPage)
		})
	}()
}

func yetToImplement(feature string) func() {
	message := fmt.Sprintf("[yellow]%s is yet to implement. Please Check in next version.", feature)
	return func() { showMessage(message) }
}

func makeButton(label string, handler func()) *tview.Button {
	btn := tview.NewButton(label).SetSelectedFunc(handler).
		SetLabelColor(tcell.ColorWhite)

	btn.SetBackgroundColor(tcell.ColorCornflowerBlue)

	return btn
}

func ignoreKeyEvt() bool {
	textInputs := []string{"*tview.InputField", "*femto.View"}
	return util.InArray(reflect.TypeOf(app.GetFocus()).String(), textInputs)
}
