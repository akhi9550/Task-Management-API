package usecase

import (
	"errors"
	"taskmanagementapi/pkg/helper"
	interfaces "taskmanagementapi/pkg/repository/interface"
	services "taskmanagementapi/pkg/usecase/interface"
	"taskmanagementapi/pkg/utils/models"

	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepository interfaces.UserRepository
}

func NewUserUseCase(repository interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepository: repository,
	}
}

func (ur *userUseCase) UserSignUp(user models.UserSignup) error {
	email, err := ur.userRepository.CheckUserExistsByEmail(user.Email)
	if email {
		return errors.New("user with this email is already exists")
	}
	if err != nil {
		return errors.New("error from check email")
	}
	hashPassword, err := helper.PasswordHash(user.Password)
	if err != nil {
		return errors.New("error in hashing password")
	}
	user.Password = hashPassword
	err = ur.userRepository.UserSignUp(user)
	if err != nil {
		return errors.New("could not add the user data")
	}
	return nil
}

func (ur *userUseCase) UserSignIn(user models.UserSignIn) (string, error) {
	exist, err := ur.userRepository.CheckUserExistsByEmail(user.Email)
	if err != nil {
		return "", err 
	}
	if !exist {
		return "", errors.New("email doesn't exist")
	}
	

	userdeatils, err := ur.userRepository.FindUserDetailsByEmail(user.Email)
	if err != nil {
		return "", errors.New("error in find user details")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userdeatils.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("password not matching")
	}
	token, err := ur.userRepository.GenerateJwtToken(userdeatils)
	if err != nil {
		return "", errors.New("couldn't create token")
	}
	return token, nil
}
