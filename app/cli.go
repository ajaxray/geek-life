package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	flag "github.com/spf13/pflag"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
	repo "github.com/ajaxray/geek-life/repository/storm"
	"github.com/ajaxray/geek-life/util"
)

var (
	app              *tview.Application
	layout, contents *tview.Flex

	statusBar         *StatusBar
	projectPane       *ProjectPane
	taskPane          *TaskPane
	taskDetailPane    *TaskDetailPane
	projectDetailPane *ProjectDetailPane

	db          *storm.DB
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository

	// Flag variables
	dbFile string
)

func init() {
	flag.StringVarP(&dbFile, "db-file", "d", "", "Specify DB file path manually.")
}

func main() {
	app = tview.NewApplication()
	flag.Parse()

	db = util.ConnectStorm(dbFile)
	defer func() {
		if err := db.Close(); err != nil {
			util.LogIfError(err, "Error in closing storm Db")
		}
	}()

	if flag.NArg() > 0 && flag.Arg(0) == "migrate" {
		migrate(db)
		fmt.Println("Database migrated successfully!")
	} else {
		projectRepo = repo.NewProjectRepository(db)
		taskRepo = repo.NewTaskRepository(db)

		layout = tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(makeTitleBar(), 2, 1, false).
			AddItem(prepareContentPages(), 0, 2, true).
			AddItem(prepareStatusBar(app), 1, 1, false)

		setKeyboardShortcuts()

		if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	}

}

func migrate(database *storm.DB) {
	util.FatalIfError(database.ReIndex(&model.Project{}), "Error in migrating Projects")
	util.FatalIfError(database.ReIndex(&model.Task{}), "Error in migrating Tasks")

	fmt.Println("Migration completed. Start geek-life normally.")
	os.Exit(0)
}

func setKeyboardShortcuts() *tview.Application {
	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvt() {
			return event
		}

		// Global shortcuts
		switch unicode.ToLower(event.Rune()) {
		case 'p':
			app.SetFocus(projectPane)
			contents.RemoveItem(taskDetailPane)
			return nil
		case 'q':
		case 't':
			app.SetFocus(taskPane)
			contents.RemoveItem(taskDetailPane)
			return nil
		}

		// Handle based on current focus. Handlers may modify event
		switch {
		case projectPane.HasFocus():
			event = projectPane.handleShortcuts(event)
		case taskPane.HasFocus():
			event = taskPane.handleShortcuts(event)
			if event != nil && projectDetailPane.isShowing() {
				event = projectDetailPane.handleShortcuts(event)
			}
		case taskDetailPane.HasFocus():
			event = taskDetailPane.handleShortcuts(event)
		}

		return event
	})
}

func prepareContentPages() *tview.Flex {
	projectPane = NewProjectPane(projectRepo)
	taskPane = NewTaskPane(projectRepo, taskRepo)
	projectDetailPane = NewProjectDetailPane()
	taskDetailPane = NewTaskDetailPane(taskRepo)

	contents = tview.NewFlex().
		AddItem(projectPane, 25, 1, true).
		AddItem(taskPane, 0, 2, false)

	return contents

}

func makeTitleBar() *tview.Flex {
	titleText := tview.NewTextView().SetText("[lime::b]Geek-life [::-]- Task Manager for geeks!").SetDynamicColors(true)
	versionInfo := tview.NewTextView().SetText("[::d]Version: 0.1.2").SetTextAlign(tview.AlignRight).SetDynamicColors(true)

	return tview.NewFlex().
		AddItem(titleText, 0, 2, false).
		AddItem(versionInfo, 0, 1, false)
}

func AskYesNo(text string, f func()) {

	activePane := app.GetFocus()
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				f()
			}
			app.SetRoot(layout, true).EnableMouse(true)
			app.SetFocus(activePane)
		})

	pages := tview.NewPages().
		AddPage("background", layout, true, true).
		AddPage("modal", modal, true, true)
	_ = app.SetRoot(pages, true).EnableMouse(true)
}
