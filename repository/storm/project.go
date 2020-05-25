package storm

import (
	"github.com/asdine/storm/v3"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
)

type projectRepository struct {
	DB *storm.DB
}

// NewProjectRepository will create an object that represent the repository.Project interface
func NewProjectRepository(db *storm.DB) repository.ProjectRepository {
	return &projectRepository{db}
}

func (repo *projectRepository) GetAll() ([]model.Project, error) {
	var projects []model.Project
	err := repo.DB.All(&projects)

	return projects, err
}

func (repo *projectRepository) GetByID(id int64) (model.Project, error) {
	return repo.getOneByField("ID", id)
}

func (repo *projectRepository) GetByTitle(title string) (model.Project, error) {
	return repo.getOneByField("Title", title)
}

func (repo *projectRepository) GetByUUID(UUID string) (model.Project, error) {
	return repo.getOneByField("CloudId", UUID)
}

func (repo *projectRepository) Create(title, UUID string) (model.Project, error) {
	project := model.Project{
		Title: title,
		UUID:  UUID,
	}

	err := repo.DB.Save(&project)
	return project, err
}

func (repo *projectRepository) Update(project *model.Project) error {
	return repo.DB.Save(project)
}

func (repo *projectRepository) Delete(project *model.Project) error {
	return repo.DB.DeleteStruct(project)
}

func (repo *projectRepository) getOneByField(fieldName string, val interface{}) (model.Project, error) {
	var project model.Project
	err := repo.DB.One(fieldName, val, &project)

	return project, err
}
