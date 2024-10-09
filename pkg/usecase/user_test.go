package usecase_test

import (
	"errors"
	"taskmanagementapi/pkg/usecase"
	"taskmanagementapi/pkg/utils/models"
	"testing"

	mockRepository "taskmanagementapi/pkg/repository/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_UserSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mockRepository.NewMockUserRepository(ctrl)
	userUseCase := usecase.NewUserUseCase(userRepo)

	testData := map[string]struct {
		input   models.UserSignup
		stub    func(*mockRepository.MockUserRepository)
		wantErr error
	}{
		"success": {
			input: models.UserSignup{
				Name:     "Akhil",
				Email:    "akhil@example.com",
				Password: "password123",
			},
			stub: func(userRepo *mockRepository.MockUserRepository) {
				userRepo.EXPECT().CheckUserExistsByEmail("akhil@example.com").Return(false, nil).Times(1)
				userRepo.EXPECT().UserSignUp(gomock.Any()).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		"user exists": {
			input: models.UserSignup{
				Name:     "Akhil",
				Email:    "akhil@example.com",
				Password: "password123",
			},
			stub: func(userRepo *mockRepository.MockUserRepository) {
				userRepo.EXPECT().CheckUserExistsByEmail("akhil@example.com").Return(true, nil).Times(1)
			},
			wantErr: errors.New("user with this email is already exists"),
		},
		"error checking user": {
			input: models.UserSignup{
				Name:     "Akhil",
				Email:    "akhil@example.com",
				Password: "password123",
			},
			stub: func(userRepo *mockRepository.MockUserRepository) {
				userRepo.EXPECT().CheckUserExistsByEmail("akhil@example.com").Return(false, errors.New("error from check email")).Times(1)
			},
			wantErr: errors.New("error from check email"),
		},
		"error adding user": {
			input: models.UserSignup{
				Name:     "Akhil",
				Email:    "akhil@example.com",
				Password: "password123",
			},
			stub: func(userRepo *mockRepository.MockUserRepository) {
				userRepo.EXPECT().CheckUserExistsByEmail("akhil@example.com").Return(false, nil).Times(1)
				userRepo.EXPECT().UserSignUp(gomock.Any()).Return(errors.New("could not add the user data")).Times(1)
			},
			wantErr: errors.New("could not add the user data"),
		},
	}

	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(userRepo)
			err := userUseCase.UserSignUp(test.input)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func TestUserSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mockRepository.NewMockUserRepository(ctrl)
	userUseCase := usecase.NewUserUseCase(mockRepo)
	mockRepo.EXPECT().CheckUserExistsByEmail("test@example.com").Return(false, nil)
	_, err := userUseCase.UserSignIn(models.UserSignIn{
		Email:    "test@example.com",
		Password: "testpassword",
	})
	assert.EqualError(t, err, "email doesn't exist")
	mockRepo.EXPECT().CheckUserExistsByEmail("test@example.com").Return(true, nil)
	mockRepo.EXPECT().FindUserDetailsByEmail("test@example.com").Return(models.UserDetails{}, errors.New("db error"))
	_, err = userUseCase.UserSignIn(models.UserSignIn{
		Email:    "test@example.com",
		Password: "testpassword",
	})
	assert.EqualError(t, err, "error in find user details")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	mockRepo.EXPECT().CheckUserExistsByEmail("test@example.com").Return(true, nil)
	mockRepo.EXPECT().FindUserDetailsByEmail("test@example.com").Return(models.UserDetails{
		Password: string(hashedPassword),
	}, nil)
	_, err = userUseCase.UserSignIn(models.UserSignIn{
		Email:    "test@example.com",
		Password: "wrongpassword", 
	})
	assert.EqualError(t, err, "password not matching")
	mockRepo.EXPECT().CheckUserExistsByEmail("test@example.com").Return(true, nil)
	mockRepo.EXPECT().FindUserDetailsByEmail("test@example.com").Return(models.UserDetails{
		Password: string(hashedPassword),
	}, nil)
	mockRepo.EXPECT().GenerateJwtToken(gomock.Any()).Return("validToken", nil)
	token, err := userUseCase.UserSignIn(models.UserSignIn{
		Email:    "test@example.com",
		Password: "correctpassword",
	})
	assert.NoError(t, err)
	assert.Equal(t, "validToken", token)
}
