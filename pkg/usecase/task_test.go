package usecase_test

import (
	"errors"
	"taskmanagementapi/pkg/usecase"
	"taskmanagementapi/pkg/utils/models"
	"testing"

	mockRepository "taskmanagementapi/pkg/repository/mock" 

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepo := mockRepository.NewMockTaskRepository(ctrl)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)

	testData := map[string]struct {
		input   models.CreateTask
		userID  string
		stub    func(*mockRepository.MockTaskRepository, models.CreateTask, string)
		wantErr error
	}{
		"success": {
			input: models.CreateTask{
				Title:       "New Task",
				Description: "Task description",
			},
			userID: "123",
			stub: func(repo *mockRepository.MockTaskRepository, task models.CreateTask, userID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1)
				repo.EXPECT().InsertTask(task, userID).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		"user does not exist": {
			input: models.CreateTask{
				Title:       "New Task",
				Description: "Task description",
			},
			userID: "123",
			stub: func(repo *mockRepository.MockTaskRepository, task models.CreateTask, userID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(false, nil).Times(1)
			},
			wantErr: errors.New("user doesn't exist"),
		},
		"repository error": {
			input: models.CreateTask{
				Title:       "New Task",
				Description: "Task description",
			},
			userID: "123",
			stub: func(repo *mockRepository.MockTaskRepository, task models.CreateTask, userID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1) 
				repo.EXPECT().InsertTask(task, userID).Return(errors.New("error from insert task")).Times(1)
			},
			wantErr: errors.New("error from insert task"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(taskRepo, test.input, test.userID)
			err := taskUseCase.CreateTask(test.input, test.userID)
			assert.Equal(t, test.wantErr, err)
		})
	}
}


func Test_GetTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepo := mockRepository.NewMockTaskRepository(ctrl)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)

	testData := map[string]struct {
		userID  string
		stub    func(*mockRepository.MockTaskRepository, string)
		want    []models.TaskDetails
		wantErr error
	}{
		"success": {
			userID: "123",
			stub: func(repo *mockRepository.MockTaskRepository, userID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1)
				repo.EXPECT().GetTasks(userID).Return([]models.TaskDetails{
					{Title: "Task1", Description: "Desc1"},
					{Title: "Task2", Description: "Desc2"},
				}, nil).Times(1)
			},
			want: []models.TaskDetails{
				{Title: "Task1", Description: "Desc1"},
				{Title: "Task2", Description: "Desc2"},
			},
			wantErr: nil,
		},
		"user does not exist": {
			userID: "123",
			stub: func(repo *mockRepository.MockTaskRepository, userID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(false, nil).Times(1)
			},
			want:    []models.TaskDetails{},
			wantErr: errors.New("user doesn't exist"),
		},
		"repository error": {
			userID: "123",
			stub: func(repo *mockRepository.MockTaskRepository, userID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1)
				repo.EXPECT().GetTasks(userID).Return([]models.TaskDetails{}, errors.New("error from get tasks")).Times(1)
			},
			want:    []models.TaskDetails{},
			wantErr: errors.New("error from get tasks"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(taskRepo, test.userID)
			tasks, err := taskUseCase.GetTasks(test.userID)
			assert.Equal(t, test.want, tasks)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func Test_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepo := mockRepository.NewMockTaskRepository(ctrl)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)

	testData := map[string]struct {
		userID  string
		taskID  string
		stub    func(*mockRepository.MockTaskRepository, string, string)
		want    models.TaskDetails
		wantErr error
	}{
		"success": {
			userID: "123",
			taskID: "456",
			stub: func(repo *mockRepository.MockTaskRepository, userID string, taskID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1)
				repo.EXPECT().CheckTaskIDExist(taskID).Return(true, nil).Times(1)
				repo.EXPECT().GetTask(userID, taskID).Return(models.TaskDetails{
					Title:       "Task1",
					Description: "Desc1",
				}, nil).Times(1)
			},
			want: models.TaskDetails{
				Title:       "Task1",
				Description: "Desc1",
			},
			wantErr: nil,
		},
		"user does not exist": {
			userID: "123",
			taskID: "456",
			stub: func(repo *mockRepository.MockTaskRepository, userID string, taskID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(false, nil).Times(1)
			},
			want:    models.TaskDetails{},
			wantErr: errors.New("user doesn't exist"),
		},
		"task does not exist": {
			userID: "123",
			taskID: "456",
			stub: func(repo *mockRepository.MockTaskRepository, userID string, taskID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1)
				repo.EXPECT().CheckTaskIDExist(taskID).Return(false, nil).Times(1)
			},
			want:    models.TaskDetails{},
			wantErr: errors.New("task doesn't exist"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(taskRepo, test.userID, test.taskID)
			task, err := taskUseCase.GetTask(test.userID, test.taskID)
			assert.Equal(t, test.want, task)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func Test_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepo := mockRepository.NewMockTaskRepository(ctrl)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)

	testData := map[string]struct {
		userID  string
		taskID  string
		input   models.CreateTask
		stub    func(*mockRepository.MockTaskRepository, string, string, models.CreateTask)
		wantErr error
	}{
		"success": {
			userID: "123",
			taskID: "456",
			input: models.CreateTask{
				Title:       "Updated Task",
				Description: "Updated Description",
			},
			stub: func(repo *mockRepository.MockTaskRepository, userID string, taskID string, task models.CreateTask) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1)
				repo.EXPECT().CheckTaskIDExist(taskID).Return(true, nil).Times(1)
				repo.EXPECT().Update(userID, taskID, task).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		"user does not exist": {
			userID: "123",
			taskID: "456",
			input: models.CreateTask{
				Title:       "Updated Task",
				Description: "Updated Description",
			},
			stub: func(repo *mockRepository.MockTaskRepository, userID string, taskID string, task models.CreateTask) {
				repo.EXPECT().CheckUserIDExist(userID).Return(false, nil).Times(1)
			},
			wantErr: errors.New("user doesn't exist"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(taskRepo, test.userID, test.taskID, test.input)
			err := taskUseCase.UpdateTask(test.userID, test.taskID, test.input)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func Test_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepo := mockRepository.NewMockTaskRepository(ctrl)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)

	testData := map[string]struct {
		userID  string
		taskID  string
		stub    func(*mockRepository.MockTaskRepository, string, string)
		wantErr error
	}{
		"success": {
			userID: "123",
			taskID: "456",
			stub: func(repo *mockRepository.MockTaskRepository, userID string, taskID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(true, nil).Times(1)
				repo.EXPECT().CheckTaskIDExist(taskID).Return(true, nil).Times(1)
				repo.EXPECT().DeleteTask(userID, taskID).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		"user does not exist": {
			userID: "123",
			taskID: "456",
			stub: func(repo *mockRepository.MockTaskRepository, userID string, taskID string) {
				repo.EXPECT().CheckUserIDExist(userID).Return(false, nil).Times(1)
			},
			wantErr: errors.New("user doesn't exist"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(taskRepo, test.userID, test.taskID)
			err := taskUseCase.DeleteTask(test.userID, test.taskID)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
