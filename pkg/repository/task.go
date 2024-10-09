package repository

import (
	"context"
	"errors"
	interfaces "taskmanagementapi/pkg/repository/interface"
	"taskmanagementapi/pkg/utils/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	TaskCollection *mongo.Collection
	UserCollection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) interfaces.TaskRepository {
	return &TaskRepository{TaskCollection: db.Collection("tasks"),
		UserCollection: db.Collection("users")}
}

func (repo *TaskRepository) CheckUserIDExist(userID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid ObjectID format")
	}
	count, err := repo.UserCollection.CountDocuments(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (tk *TaskRepository) InsertTask(task models.CreateTask, userID string) error {
	currentTime := time.Now().Format(time.RFC3339)
	newTask := bson.M{
		"user_id":     userID,
		"title":       task.Title,
		"description": task.Description,
		"created_at":  currentTime,
	}
	_, err := tk.TaskCollection.InsertOne(context.TODO(), newTask)
	if err != nil {
		return err
	}
	return nil
}

func (tk *TaskRepository) GetTasks(userID string) ([]models.TaskDetails, error) {
	cursor, err := tk.TaskCollection.Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []models.TaskDetails
	for cursor.Next(context.TODO()) {
		var task models.TaskDetails
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, cursor.Err()
}

func (tk *TaskRepository) CheckTaskIDExist(taskID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return false, errors.New("invalid ObjectID format")
	}
	count, err := tk.TaskCollection.CountDocuments(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (tk *TaskRepository) GetTask(userID, taskID string) (models.TaskDetails, error) {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.TaskDetails{}, err
	}
	filter := bson.M{
		"user_id": userID,
		"_id":     objID,
	}
	var task models.TaskDetails

	err = tk.TaskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.TaskDetails{}, nil
		}
		return models.TaskDetails{}, err
	}
	return task, nil
}

func (tk *TaskRepository) Update(userID, taskID string, task models.CreateTask) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"user_id": userID,
		"_id":     objID,
	}

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
		},
	}
	_, err = tk.TaskCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (tk *TaskRepository) DeleteTask(userID, taskID string) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"user_id": userID,
		"_id":     objID,
	}
	result, err := tk.TaskCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return err
	}
	return nil
}
