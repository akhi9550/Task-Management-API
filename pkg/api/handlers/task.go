package handlers

import (
	services "taskmanagementapi/pkg/usecase/interface"
	"taskmanagementapi/pkg/utils/models"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	TaskUseCase services.TaskUseCase
}

func NewTaskHandler(useCase services.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		TaskUseCase: useCase,
	}
}

func (tk *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var task models.CreateTask
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := validator.New().Struct(task)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Constraints not statisfied"})
	}
	userID := c.Locals("user_id").(string)
	err = tk.TaskUseCase.CreateTask(task, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Task creation failed", "message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Task created"})
}

func (tk *TaskHandler) GetTasks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	tasks, err := tk.TaskUseCase.GetTasks(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Tasks Retrieve failed", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully Get Tasks", "data": tasks})
}

func (tk *TaskHandler) GetTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	userID := c.Locals("user_id").(string)
	task, err := tk.TaskUseCase.GetTask(userID, taskID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Task Retrieve failed", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully Get Task", "data": task})
}

func (tk *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	userID := c.Locals("user_id").(string)
	var task models.CreateTask
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := validator.New().Struct(task)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Constraints not statisfied"})
	}
	err = tk.TaskUseCase.UpdateTask(userID, taskID, task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Task Update failed", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully Update Task"})
}

func (tk *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	userID := c.Locals("user_id").(string)
	err := tk.TaskUseCase.DeleteTask(userID, taskID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Task Delete failed", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully Deleted Task"})
}
