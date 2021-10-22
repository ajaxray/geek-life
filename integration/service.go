package integration

import (
	"github.com/ajaxray/geek-life/integration/gtasks"
	repo "github.com/ajaxray/geek-life/repository"
)

// TaskService interface defines methods of a Task Service
type TaskService interface {
	// Key returns a string, the identifier keyword of the service
	Key() string

	// Integrate prepares auth, data-source etc and returns service data for storing
	Integrate() (map[string]interface{}, error)

	// Sync usages the integration data and syncs Projects and Tasks
	Sync(repo.ProjectRepository, repo.TaskRepository, map[string]interface{}) error

	// Clear removes integration data. Can be used for start over
	// Can be used for cleanup external data, any files created etc.
	Clear() error
}

func FindService(key string) TaskService {
	switch key {
	case gtasks.IntegrationKey:
		return gtasks.GTasks{}
	default:
		panic("No integration was defined for " + key)
	}
}
