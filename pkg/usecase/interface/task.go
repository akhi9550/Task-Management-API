package interfaces

import "taskmanagementapi/pkg/utils/models"

type TaskUseCase interface {
	CreateTask(models.CreateTask, string) error
	GetTasks(string) ([]models.TaskDetails, error)
	GetTask(string, string) (models.TaskDetails, error)
	UpdateTask(string,string,models.CreateTask)error
	DeleteTask(string,string)error
}
