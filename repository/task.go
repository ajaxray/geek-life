package repository

import (
	"time"

	"github.com/ajaxray/geek-life/model"
)

type TaskRepository interface {
	GetAll() ([]model.Task, error)
	GetAllByProject(project model.Project) ([]model.Task, error)
	GetAllByDate(from, to time.Time) ([]model.Task, error)
	GetByID(ID string) (model.Task, error)
	GetByUUID(UUID string) (model.Task, error)
	Create(project model.Project, title, details, UUID string, dueDate int64) (model.Task, error)
	Update(t *model.Task) error
	UpdateField(t *model.Task, field string, value interface{}) error
	Delete(t *model.Task) error
}
