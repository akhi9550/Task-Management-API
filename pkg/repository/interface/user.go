package interfaces

import "taskmanagementapi/pkg/utils/models"

type UserRepository interface {
	CheckUserExistsByEmail(string) (bool, error)
	UserSignUp(models.UserSignup) error
	FindUserDetailsByEmail(string) (models.UserDetails, error)
	GenerateJwtToken(user models.UserDetails) (string, error)
}
