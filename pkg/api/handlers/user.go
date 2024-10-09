package handlers

import (
	services "taskmanagementapi/pkg/usecase/interface"
	"taskmanagementapi/pkg/utils/models"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserUseCase services.UserUseCase
}

func NewUserHandler(useCase services.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: useCase,
	}
}

func (ur *UserHandler) UserSignUp(c *fiber.Ctx) error {
	var user models.UserSignup
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := validator.New().Struct(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Constraints not statisfied"})
	}
	err = ur.UserUseCase.UserSignUp(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User creation failed", "message": err.Error()})

	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created"})
}

func (ur *UserHandler) UserSignIn(c *fiber.Ctx) error {
	var user models.UserSignIn
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := validator.New().Struct(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Constraints not statisfied"})
	}
	token, err := ur.UserUseCase.UserSignIn(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User signIn failed", "message": err.Error()})

	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User signIn Successful", "token": token})
}

func (u *UserHandler) UserSignOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), 
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User signout successful"})
}
