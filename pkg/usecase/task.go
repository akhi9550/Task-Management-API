package usecase

import (
	"errors"
	interfaces "taskmanagementapi/pkg/repository/interface"
	services "taskmanagementapi/pkg/usecase/interface"
	"taskmanagementapi/pkg/utils/models"
)

type TaskUseCase struct {
	taskRepository interfaces.TaskRepository
}

func NewTaskUseCase(repository interfaces.TaskRepository) services.TaskUseCase {
	return &TaskUseCase{
		taskRepository: repository,
	}
}

func (tk *TaskUseCase) CreateTask(task models.CreateTask, userID string) error {
	exist, err := tk.taskRepository.CheckUserIDExist(userID)
	if !exist {
		return errors.New("user doesn't exist")
	}
	if err != nil {
		return err
	}
	err = tk.taskRepository.InsertTask(task, userID)
	if err != nil {
		return errors.New("error from insert task")
	}
	return nil
}

func (tk *TaskUseCase) GetTasks(userID string) ([]models.TaskDetails, error) {
	exist, err := tk.taskRepository.CheckUserIDExist(userID)
	if !exist {
		return []models.TaskDetails{}, errors.New("user doesn't exist")
	}
	if err != nil {
		return []models.TaskDetails{}, err
	}
	tasks, err := tk.taskRepository.GetTasks(userID)
	if err != nil {
		return []models.TaskDetails{}, errors.New("error from get tasks")
	}
	return tasks, nil
}

func (tk *TaskUseCase) GetTask(userID, taskID string) (models.TaskDetails, error) {
	existUserID, err := tk.taskRepository.CheckUserIDExist(userID)
	if !existUserID {
		return models.TaskDetails{}, errors.New("user doesn't exist")
	}
	if err != nil {
		return models.TaskDetails{}, err
	}
	existTaskID, err := tk.taskRepository.CheckTaskIDExist(taskID)
	if !existTaskID {
		return models.TaskDetails{}, errors.New("task doesn't exist")
	}
	if err != nil {
		return models.TaskDetails{}, err
	}
	task, err := tk.taskRepository.GetTask(userID, taskID)
	if err != nil {
		return models.TaskDetails{}, errors.New("error from get task")
	}
	return task, nil
}

func (tk *TaskUseCase) UpdateTask(userID, taskID string, task models.CreateTask) error {
	existUserID, err := tk.taskRepository.CheckUserIDExist(userID)
	if !existUserID {
		return errors.New("user doesn't exist")
	}
	if err != nil {
		return err
	}
	existTaskID, err := tk.taskRepository.CheckTaskIDExist(taskID)
	if !existTaskID {
		return errors.New("task doesn't exist")
	}
	if err != nil {
		return err
	}
	err = tk.taskRepository.Update(userID, taskID, task)
	if err != nil {
		return errors.New("error from update title")
	}
	return nil
}

func (tk *TaskUseCase) DeleteTask(userID, taskID string) error {
	existUserID, err := tk.taskRepository.CheckUserIDExist(userID)
	if !existUserID {
		return errors.New("user doesn't exist")
	}
	if err != nil {
		return err
	}
	existTaskID, err := tk.taskRepository.CheckTaskIDExist(taskID)
	if !existTaskID {
		return errors.New("task doesn't exist")
	}
	if err != nil {
		return err
	}
	err = tk.taskRepository.DeleteTask(userID, taskID)
	if err != nil {
		return errors.New("error from delete task")
	}
	return nil
}
