package main

import (
	"fmt"
	"os"

	"github.com/ajaxray/geek-life/integration"
	repo "github.com/ajaxray/geek-life/repository/storm"
	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	flag "github.com/spf13/pflag"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
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
	dbFile           string
	sync             bool
	clearIntegration bool
)

func init() {
	flag.StringVarP(&dbFile, "db-file", "d", "", "Specify DB file path manually.")
	flag.BoolVar(&sync, "sync", false, "Sync with defined integration service then start")
	flag.BoolVar(&clearIntegration, "clear", false, "Clear a defined integration")
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

	projectRepo = repo.NewProjectRepository(db)
	taskRepo = repo.NewTaskRepository(db)

	if flag.NArg() > 0 && flag.Arg(0) == "migrate" {
		migrate(db)
		fmt.Println("Database migrated successfully!")
	} else if flag.NArg() == 2 && flag.Arg(0) == "integrate" {
		integration.HandleIntegration(db, flag.Arg(1), clearIntegration)
	} else {
		if sync {
			integration.HandleSync(projectRepo, taskRepo, db)
			os.Exit(0)
		}

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
		switch event.Rune() {
		case 'p':
			app.SetFocus(projectPane)
			return nil
		case 't':
			app.SetFocus(taskPane)
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
