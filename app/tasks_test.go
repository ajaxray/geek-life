package main

import (
	"testing"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
	"github.com/rivo/tview"
)

type deleteRecordingTaskRepository struct {
	repository.TaskRepository
	deletedIDs []int64
}

func (repo *deleteRecordingTaskRepository) Delete(task *model.Task) error {
	repo.deletedIDs = append(repo.deletedIDs, task.ID)
	return nil
}

func TestClearCompletedTasksKeepsTaskListAligned(t *testing.T) {
	previousStatusBar := statusBar
	statusBar = &StatusBar{}
	t.Cleanup(func() { statusBar = previousStatusBar })

	tasks := []model.Task{
		{ID: 1, Title: "completed one", Completed: true},
		{ID: 2, Title: "pending one"},
		{ID: 3, Title: "completed two", Completed: true},
		{ID: 4, Title: "pending two"},
	}
	list := tview.NewList().ShowSecondaryText(false)
	for _, task := range tasks {
		list.AddItem(task.Title, "", 0, nil)
	}
	repo := &deleteRecordingTaskRepository{}
	pane := &TaskPane{list: list, tasks: tasks, taskRepo: repo}
	want := []model.Task{
		{ID: 2, Title: "pending one"},
		{ID: 4, Title: "pending two"},
	}

	pane.ClearCompletedTasks()

	deleted := map[int64]bool{}
	for _, id := range repo.deletedIDs {
		deleted[id] = true
	}
	if len(repo.deletedIDs) != 2 || !deleted[1] || !deleted[3] {
		t.Errorf("deleted task IDs = %v, want tasks 1 and 3", repo.deletedIDs)
	}
	if len(pane.tasks) != len(want) {
		t.Errorf("remaining task count = %d, want %d", len(pane.tasks), len(want))
	}
	if list.GetItemCount() != len(want) {
		t.Errorf("visible task count = %d, want %d", list.GetItemCount(), len(want))
	}
	for i := range want {
		if i < len(pane.tasks) && pane.tasks[i].ID != want[i].ID {
			t.Errorf("pane.tasks[%d].ID = %d, want %d", i, pane.tasks[i].ID, want[i].ID)
		}
		if i < list.GetItemCount() {
			title, _ := list.GetItemText(i)
			if title != want[i].Title {
				t.Errorf("visible item %d = %q, want %q", i, title, want[i].Title)
			}
		}
	}
}
