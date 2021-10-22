package gtasks

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ajaxray/geek-life/model"
	repo "github.com/ajaxray/geek-life/repository"
	"github.com/ajaxray/geek-life/util"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gTaskAPI "google.golang.org/api/tasks/v1"
)

// integrationKey will be used as prefix for all integration affairs
const IntegrationKey = "gtasks"

// Create from https://console.cloud.google.com/apis/credentials/oauthclient/
// Download as JSON and put as credentials.json in integration/gtasks directory
//go:embed credentials.json
var credentials []byte
var taskCompleted = map[string]bool{"completed": true, "needsAction": false}

type GTasks struct{}

func (GT GTasks) Integrate() (map[string]interface{}, error) {
	serviceData := make(map[string]interface{})

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials, gTaskAPI.TasksScope)
	if err != nil {
		serviceData["error"] = "Unable to parse client secret to config"
		return serviceData, err
	}

	token, err := json.Marshal(getTokenFromWeb(config))
	util.FatalIfError(err, "Couldn't receive access token properly!")
	serviceData["token"] = string(token)

	return serviceData, nil
}

func (GT GTasks) Sync(projectRepo repo.ProjectRepository, taskRepo repo.TaskRepository, serviceData map[string]interface{}) error {
	service := GT.getService(serviceData)

	lists, err := service.Tasklists.List().Do()
	util.FatalIfError(err, "Unable to retrieve task lists.")

	fmt.Println("Syncing Task Lists:")
	for _, list := range lists.Items {
		fmt.Printf("# %s (%s)\n", list.Title, list.Id)
		var project model.Project
		var projectErr error

		project, projectErr = projectRepo.GetByIntegrationID(list.Id)

		if projectErr != nil {
			fmt.Printf("# Creating Project: %s \n", list.Title)
			project = GT.createProject(projectRepo, list)
		} else {
			fmt.Printf("# Project Exixts: %s (Updated: %s) \n", list.Title, list.Updated)
			// @TODO: Check for Update
		}

		GT.syncProjectTasks(project, service, taskRepo)
	}

	return nil
}

func (GT GTasks) syncProjectTasks(project model.Project, service *gTaskAPI.Service, taskRepo repo.TaskRepository) {
	tasks, err := service.Tasks.List(project.IntegrationID).Do()
	util.FatalIfError(err, "Unable to retrieve tasks of list")

	for _, task := range tasks.Items {
		// fmt.Printf("### %s (%s)\n", task.Title, task.Id)
		localTask, err := taskRepo.GetByIntegrationID(task.Id)
		if err == nil {
			fmt.Printf("# Task Exixts: %s (Updated: %s) \n", localTask.Title, time.Unix(localTask.ModifiedAt, 0))

			if task.Deleted {
				_ = taskRepo.Delete(&localTask)
				continue
			}

			GT.syncTaskUpdate(task, localTask, project, taskRepo, service)
		} else {
			if taskCompleted[task.Status] || task.Deleted {
				fmt.Printf("# Skipping completed or deleted task: %s (Updated: %s) \n", task.Title, task.Updated)
				continue
			}

			localTask, err := taskRepo.Create(project, task.Title, task.Notes, task.Id, gTaskToUnixTime(task.Due))
			util.LogIfError(err, "Failed to create task: %s", task.Title)

			err = taskRepo.UpdateField(&localTask, "ModifiedAt", gTaskToUnixTime(task.Updated))
			util.LogIfError(err, "Failed to update modification time: %s", task.Title)

			fmt.Printf("# Task Created: %s \n", localTask.Title)
		}
	}
}

func (GT GTasks) syncTaskUpdate(remoteTask *gTaskAPI.Task, localTask model.Task, project model.Project, taskRepo repo.TaskRepository, service *gTaskAPI.Service) {

	fmt.Printf("Comparing - %s (Local %s <=> Remote %s) \n", localTask.Title, unixToGtaskTime(localTask.ModifiedAt), remoteTask.Updated)

	if gTaskToUnixTime(remoteTask.Updated) < localTask.ModifiedAt {
		fmt.Println("Modified Locally")
		GT.updateRemoteByLocalTask(remoteTask, localTask, service, project)
	} else if gTaskToUnixTime(remoteTask.Updated) > localTask.ModifiedAt {
		fmt.Println("Modified Remotely")
		GT.updateLocalByRemoteTask(localTask, remoteTask, taskRepo)
	}
}

func (GT GTasks) updateLocalByRemoteTask(localTask model.Task, remoteTask *gTaskAPI.Task, taskRepo repo.TaskRepository) {
	localTask.ModifiedAt = gTaskToUnixTime(remoteTask.Updated)
	localTask.Title = remoteTask.Title
	localTask.Details = remoteTask.Notes
	localTask.Completed = remoteTask.Status == "completed"
	if remoteTask.Due != "" {
		localTask.DueDate = gTaskToUnixTime(remoteTask.Due)
	} else {
		localTask.DueDate = 0
	}

	err := taskRepo.Update(&localTask)
	util.LogIfError(err, "Error in saving updated Task")
}

func (GT GTasks) updateRemoteByLocalTask(remoteTask *gTaskAPI.Task, localTask model.Task, service *gTaskAPI.Service, project model.Project) {
	remoteTask.Updated = unixToGtaskTime(localTask.ModifiedAt)
	remoteTask.Title = localTask.Title
	remoteTask.Notes = localTask.Details
	remoteTask.Status = "needsAction"
	if localTask.Completed {
		remoteTask.Status = "completed"
	}

	if localTask.DueDate != 0 {
		remoteTask.Due = unixToGtaskTime(localTask.DueDate)
	} else {
		remoteTask.Due = ""
	}

	_, err := service.Tasks.Update(project.IntegrationID, remoteTask.Id, remoteTask).Do()
	util.LogIfError(err, "Failed to update remote task")
}

func (GT GTasks) getService(serviceData map[string]interface{}) *gTaskAPI.Service {
	config, err := google.ConfigFromJSON([]byte(credentials), gTaskAPI.TasksScope)

	token := &oauth2.Token{}
	tokenData, _ := serviceData["token"]
	err = json.Unmarshal([]byte(tokenData.(string)), &token)
	util.FatalIfError(err, "Could not use access token. Try running with 'integrate gtasks' again.")

	client := config.Client(context.Background(), token)

	service, err := gTaskAPI.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gTaskAPI Client %v", err)
	}

	return service
}

func (GT GTasks) Clear() error {
	return nil
}

func (GT GTasks) Key() string {
	return IntegrationKey
}

func (GT GTasks) createProject(projectRepo repo.ProjectRepository, list *gTaskAPI.TaskList) model.Project {
	newProject, _ := projectRepo.Create(list.Title, list.Id)
	_ = projectRepo.UpdateField(&newProject, "ModifiedAt", gTaskToUnixTime(list.Updated))

	return newProject
}
