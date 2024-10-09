package interfaces

import "taskmanagementapi/pkg/utils/models"

type TaskRepository interface {
	CheckUserIDExist(string) (bool, error)
	InsertTask(models.CreateTask, string) error
	GetTasks(string) ([]models.TaskDetails, error)
	CheckTaskIDExist(string) (bool, error)
	GetTask(string, string) (models.TaskDetails, error)
	Update(string, string, models.CreateTask) error
	DeleteTask(string, string) error
}
