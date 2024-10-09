package repository_test

import (
	"taskmanagementapi/pkg/repository"
	"taskmanagementapi/pkg/utils/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCheckUserIDExist(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("valid ObjectID exists", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
				{Key: "n", Value: 1},
				{Key: "_id", Value: userID},
			}),
		)

		userRepo := repository.NewTaskRepository(mt.Client.Database("test"))
		exists, err := userRepo.CheckUserIDExist(userID.Hex())

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	mt.Run("valid ObjectID does not exist", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
				{Key: "n", Value: 0},
			}),
		)
		userRepo := repository.NewTaskRepository(mt.Client.Database("test"))
		exists, err := userRepo.CheckUserIDExist(userID.Hex())
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	mt.Run("invalid ObjectID format", func(mt *mtest.T) {
		userRepo := repository.NewTaskRepository(mt.Client.Database("test"))
		_, err := userRepo.CheckUserIDExist("invalid_id")

		assert.Error(t, err)
		assert.Equal(t, "invalid ObjectID format", err.Error())
	})
}

func TestInsertTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("successfully insert task", func(mt *mtest.T) {
		task := models.CreateTask{
			Title:       "Test Task",
			Description: "This is a test task",
		}
		userID := primitive.NewObjectID().Hex()
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		tk := repository.NewTaskRepository(mt.Client.Database("test"))

		err := tk.InsertTask(task, userID)

		assert.NoError(t, err)
	})

	mt.Run("error during InsertOne operation", func(mt *mtest.T) {
		task := models.CreateTask{
			Title:       "Test Task",
			Description: "This is a test task",
		}
		userID := primitive.NewObjectID().Hex()
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate key error",
		}))

		tk := repository.NewTaskRepository(mt.Client.Database("test"))

		err := tk.InsertTask(task, userID)

		assert.EqualError(t, err, "duplicate key error")
	})
}

func TestGetTasks(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("valid user ID with tasks", func(mt *mtest.T) {
		userID := "valid_user_id"
		taskID1 := primitive.NewObjectID()
		taskID2 := primitive.NewObjectID()
		now := time.Now().UTC().Truncate(time.Second)

		task1 := models.TaskDetails{
			ID:          taskID1.Hex(),
			Title:       "Test Task 1",
			Description: "Test Description 1",
			CreatedAt:   now,
		}
		task2 := models.TaskDetails{
			ID:          taskID2.Hex(),
			Title:       "Test Task 2",
			Description: "Test Description 2",
			CreatedAt:   now,
		}
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.tasks", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: taskID1},
				{Key: "user_id", Value: userID},
				{Key: "title", Value: task1.Title},
				{Key: "description", Value: task1.Description},
				{Key: "created_at", Value: task1.CreatedAt},
			}),
			mtest.CreateCursorResponse(1, "test.tasks", mtest.NextBatch, bson.D{
				{Key: "_id", Value: taskID2},
				{Key: "user_id", Value: userID},
				{Key: "title", Value: task2.Title},
				{Key: "description", Value: task2.Description},
				{Key: "created_at", Value: task2.CreatedAt},
			}),
			mtest.CreateCursorResponse(0, "test.tasks", mtest.NextBatch),
		)

		taskRepo := repository.NewTaskRepository(mt.Client.Database("test"))
		tasks, err := taskRepo.GetTasks(userID)
		for i := range tasks {
			tasks[i].CreatedAt = tasks[i].CreatedAt.UTC().Truncate(time.Second)
		}
		task1.CreatedAt = task1.CreatedAt.Truncate(time.Second)
		task2.CreatedAt = task2.CreatedAt.Truncate(time.Second)

		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
		assert.Equal(t, task1, tasks[0])
		assert.Equal(t, task2, tasks[1])
	})

	mt.Run("valid user ID without tasks", func(mt *mtest.T) {
		userID := "valid_user_id"
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.tasks", mtest.FirstBatch),
			mtest.CreateCursorResponse(0, "test.tasks", mtest.NextBatch),
		)

		taskRepo := repository.NewTaskRepository(mt.Client.Database("test"))
		tasks, err := taskRepo.GetTasks(userID)

		assert.NoError(t, err)
		assert.Len(t, tasks, 0)
	})

	mt.Run("invalid user ID format", func(mt *mtest.T) {
		invalidUserID := "invalid_user_id"

		taskRepo := repository.NewTaskRepository(mt.Client.Database("test"))
		tasks, err := taskRepo.GetTasks(invalidUserID)

		assert.Error(t, err)
		assert.Nil(t, tasks)
	})
}

func TestCheckTaskIDExist(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("valid ObjectID exists", func(mt *mtest.T) {
		taskID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.tasks", mtest.FirstBatch, bson.D{
				{Key: "n", Value: 1},
				{Key: "_id", Value: taskID},
			}),
		)

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		exists, err := tk.CheckTaskIDExist(taskID.Hex())

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	mt.Run("valid ObjectID does not exist", func(mt *mtest.T) {
		taskID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.tasks", mtest.FirstBatch, bson.D{
				{Key: "n", Value: 0},
			}),
		)

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		exists, err := tk.CheckTaskIDExist(taskID.Hex())

		assert.NoError(t, err)
		assert.False(t, exists)
	})

	mt.Run("invalid ObjectID format", func(mt *mtest.T) {
		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		_, err := tk.CheckTaskIDExist("invalid_id")
		assert.Error(t, err)
		assert.Equal(t, "invalid ObjectID format", err.Error())
	})
}

func TestGetTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful retrieval of task", func(mt *mtest.T) {
		userID := "6705824f80a09eb0313f0e42"
		taskID := primitive.NewObjectID().Hex()
		task := models.TaskDetails{
			ID:          taskID,
			Title:       "Task 1",
			Description: "Description for Task 1",
			CreatedAt:   time.Now(),
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.tasks", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: task.ID},
			{Key: "title", Value: task.Title},
			{Key: "description", Value: task.Description},
			{Key: "created_at", Value: task.CreatedAt},
		}))

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		retrievedTask, err := tk.GetTask(userID, task.ID)

		assert.NoError(t, err)
		assert.Equal(t, task.ID, retrievedTask.ID)
		assert.Equal(t, task.Title, retrievedTask.Title)
	})

	mt.Run("task not found", func(mt *mtest.T) {
		userID := "6705824f80a09eb0313f0e42"
		taskID := primitive.NewObjectID().Hex()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.tasks", mtest.FirstBatch))

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		retrievedTask, err := tk.GetTask(userID, taskID)

		assert.NoError(t, err)
		assert.Empty(t, retrievedTask)
	})

	mt.Run("invalid ObjectID format", func(mt *mtest.T) {
		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		_, err := tk.GetTask("6705824f80a09eb0313f0e42", "invalid_id")

		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully update task", func(mt *mtest.T) {
		taskID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		task := models.CreateTask{
			Title:       "Updated Task",
			Description: "Updated description",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse()) 

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		err := tk.Update(userID, taskID, task)

		assert.NoError(t, err)
	})

	mt.Run("error during update operation", func(mt *mtest.T) {
		taskID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()
		task := models.CreateTask{
			Title:       "Task",
			Description: "Description",
		}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate key error",
		}))

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		err := tk.Update(userID, taskID, task)

		assert.Error(t, err)
		assert.EqualError(t, err, "duplicate key error")
	})

	mt.Run("invalid ObjectID format", func(mt *mtest.T) {
		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		err := tk.Update("userID", "invalid_id", models.CreateTask{})
		assert.Error(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully delete task", func(mt *mtest.T) {
		taskID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		err := tk.DeleteTask(userID, taskID)

		assert.NoError(t, err)
	})

	mt.Run("task not found", func(mt *mtest.T) {
		taskID := primitive.NewObjectID().Hex()
		userID := primitive.NewObjectID().Hex()

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		err := tk.DeleteTask(userID, taskID)

		assert.NoError(t, err)
	})

	mt.Run("invalid ObjectID format", func(mt *mtest.T) {
		tk := repository.NewTaskRepository(mt.Client.Database("test"))
		err := tk.DeleteTask("6705824f80a09eb0313f0e42", "invalid_id")
		assert.Error(t, err)
	})
}
