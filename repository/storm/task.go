package storm

import (
	"time"

	"github.com/asdine/storm/v3"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
)

type taskRepository struct {
	DB *storm.DB
}

// NewTaskRepository will create an object that represent the repository.Task interface
func NewTaskRepository(db *storm.DB) repository.TaskRepository {
	return &taskRepository{db}
}

func (t *taskRepository) GetAll() ([]model.Task, error) {
	panic("implement me")
}

func (t *taskRepository) GetAllByProject(project model.Project) ([]model.Task, error) {
	var tasks []model.Task
	//err = db.Find("ProjetID", project.ID, &tasks, storm.Limit(10), storm.Skip(10), storm.Reverse())
	err := t.DB.Find("ProjectID", project.ID, &tasks)

	return tasks, err
}

func (t *taskRepository) GetAllByDate(from, to time.Time) ([]model.Task, error) {
	panic("implement me")
}

func (t *taskRepository) GetByID(ID string) (model.Task, error) {
	panic("implement me")
}

func (t *taskRepository) GetByUUID(UUID string) (model.Task, error) {
	panic("implement me")
}

func (t *taskRepository) Create(project model.Project, title, details, UUID string, dueDate int64) (model.Task, error) {
	task := model.Task{
		ProjectID: project.ID,
		Title:     title,
		Details:   details,
		UUID:      UUID,
		DueDate:   dueDate,
	}

	err := t.DB.Save(&task)
	return task, err
}

func (t *taskRepository) Update(task *model.Task) error {
	return t.DB.Update(task)
}

func (t *taskRepository) UpdateField(task *model.Task, field string, value interface{}) error {
	return t.DB.UpdateField(task, field, value)
}

func (t *taskRepository) Delete(task *model.Task) error {
	return t.DB.DeleteStruct(task)
}
