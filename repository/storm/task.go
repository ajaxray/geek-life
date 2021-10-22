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

func (t *taskRepository) GetAllByDate(date time.Time) ([]model.Task, error) {
	var tasks []model.Task

	if date.IsZero() {
		var allTasks []model.Task
		err := t.DB.AllByIndex("ProjectID", &allTasks)
		for _, t := range allTasks {
			if t.DueDate == 0 {
				tasks = append(tasks, t)
			}
		}

		return tasks, err
	} else {
		err := t.DB.Find("DueDate", getRoundedDueDate(date), &tasks)
		return tasks, err
	}
}

func (t *taskRepository) GetAllByDateRange(from, to time.Time) ([]model.Task, error) {
	var tasks []model.Task

	err := t.DB.Range("DueDate", getRoundedDueDate(from), getRoundedDueDate(to), &tasks)
	return tasks, err
}

func (t *taskRepository) GetByID(ID int64) (model.Task, error) {
	return t.getOneByField("ID", ID)
}

func (t *taskRepository) GetByIntegrationID(integrationID string) (model.Task, error) {
	return t.getOneByField("IntegrationID", integrationID)
}

func (t *taskRepository) Create(project model.Project, title, details, integrationID string, dueDate int64) (model.Task, error) {
	task := model.Task{
		ProjectID:     project.ID,
		Title:         title,
		Details:       details,
		IntegrationID: integrationID,
		DueDate:       dueDate,
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

func getRoundedDueDate(date time.Time) int64 {
	if date.IsZero() {
		return 0
	}

	return date.Unix()
}

func (repo *taskRepository) getOneByField(fieldName string, val interface{}) (model.Task, error) {
	var task model.Task
	err := repo.DB.One(fieldName, val, &task)

	return task, err
}
