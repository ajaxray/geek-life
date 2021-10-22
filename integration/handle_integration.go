package integration

import (
	"fmt"

	repo "github.com/ajaxray/geek-life/repository"
	"github.com/ajaxray/geek-life/util"
	"github.com/asdine/storm/v3"
)

// HandleIntegration finds the service and manage integration
func HandleIntegration(db *storm.DB, serviceKey string, isCleanup bool) {
	service := FindService(serviceKey)
	if isCleanup {
		cleanup(db, service)
	} else {
		data, err := service.Integrate()
		util.FatalIfError(err, "Failed to establish integration!")

		db.Set("config", "integrated", serviceKey)
		db.Set("config", serviceKey+"-data", data)
		fmt.Println("Integration completed. Now run app with --sync option.")
	}
}

func cleanup(db *storm.DB, service TaskService) {
	db.Delete("config", "integrated")
	db.Delete("config", service.Key()+"-data")
	service.Clear()
	fmt.Println("Cleared integration with " + service.Key())
}

// HandleSync finds the service and manage integration
func HandleSync(projectRepo repo.ProjectRepository, taskRepo repo.TaskRepository, db *storm.DB) {
	var serviceKey string
	err := db.Get("config", "integrated", &serviceKey)
	util.FatalIfError(err, "No service was integrated to sync with")

	fmt.Println(serviceKey + " Found as integrated service.")

	var serviceData map[string]interface{}
	db.Get("config", serviceKey+"-data", &serviceData)

	FindService(serviceKey).Sync(projectRepo, taskRepo, serviceData)
}
