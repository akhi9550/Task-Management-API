package interfaces

import "taskmanagementapi/pkg/utils/models"

type UserUseCase interface {
	UserSignUp(models.UserSignup) error
	UserSignIn(models.UserSignIn) (string, error)
}
