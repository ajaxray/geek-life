package repository

import "github.com/ajaxray/geek-life/model"

type ProjectRepository interface {
	GetAll() ([]model.Project, error)
	GetByID(id int64) (model.Project, error)
	GetByTitle(title string) (model.Project, error)
	GetByUUID(UUID string) (model.Project, error)
	Create(title, UUID string) (model.Project, error)
	Update(p *model.Project) error
	UpdateField(p *model.Project, field string, value interface{}) error
	Delete(p *model.Project) error
}
