package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskmanagementapi/pkg/api/handlers"
	"taskmanagementapi/pkg/usecase/mock"
	"taskmanagementapi/pkg/utils/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_UserSignUp(t *testing.T) {
	testCases := map[string]struct {
		input         models.UserSignup
		buildStub     func(useCaseMock *mock.MockUserUseCase, user models.UserSignup)
		checkResponse func(t *testing.T, resp *http.Response)
	}{
		"Valid User SignUp": {
			input: models.UserSignup{
				Name:     "Arun C",
				Email:    "arun@gmail.com",
				Password: "password123",
			},
			buildStub: func(useCaseMock *mock.MockUserUseCase, user models.UserSignup) {
				useCaseMock.EXPECT().UserSignUp(gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			},
		},
		"Invalid User Input": {
			input: models.UserSignup{
				Name:     "",
				Email:    "arun@gmail.com",
				Password: "password123",
			},
			buildStub: func(useCaseMock *mock.MockUserUseCase, user models.UserSignup) {
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			},
		},
		"User Creation Failure": {
			input: models.UserSignup{
				Name:     "Arun C",
				Email:    "arun@gmail.com",
				Password: "password123",
			},
			buildStub: func(useCaseMock *mock.MockUserUseCase, user models.UserSignup) {
				useCaseMock.EXPECT().UserSignUp(gomock.Any()).Times(1).Return(errors.New("user creation failed"))
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			},
		},
	}

	for testName, test := range testCases {
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUseCase := mock.NewMockUserUseCase(ctrl)
			test.buildStub(mockUseCase, test.input)
			userHandler := handlers.NewUserHandler(mockUseCase)
			app := fiber.New()
			app.Post("/signup", userHandler.UserSignUp)
			jsonData, err := json.Marshal(test.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)
			req := httptest.NewRequest("POST", "/signup", body)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			test.checkResponse(t, resp)
		})
	}
}

func Test_UserSignIn(t *testing.T) {
    testCases := map[string]struct {
        input         models.UserSignIn
        buildStub     func(useCaseMock *mock.MockUserUseCase, user models.UserSignIn)
        checkResponse func(t *testing.T, resp *http.Response)
    }{
        "Valid User SignIn": {
            input: models.UserSignIn{
                Email:    "arun@gmail.com",
                Password: "password123",
            },
            buildStub: func(useCaseMock *mock.MockUserUseCase, user models.UserSignIn) {
                useCaseMock.EXPECT().UserSignIn(user).Times(1).Return("mocked_jwt_token", nil)
            },
            checkResponse: func(t *testing.T, resp *http.Response) {
                assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
        
            },
        },
        "Invalid User Input": {
            input: models.UserSignIn{
                Email:    "",
                Password: "password123",
            },
            buildStub: func(useCaseMock *mock.MockUserUseCase, user models.UserSignIn) {
            },
            checkResponse: func(t *testing.T, resp *http.Response) {
                assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
            },
        },
        "User SignIn Failure": {
            input: models.UserSignIn{
                Email:    "arun@gmail.com",
                Password: "wrongpassword",
            },
            buildStub: func(useCaseMock *mock.MockUserUseCase, user models.UserSignIn) {
                useCaseMock.EXPECT().UserSignIn(user).Times(1).Return("", errors.New("user signIn failed"))
            },
            checkResponse: func(t *testing.T, resp *http.Response) {
                assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
              
            },
        },
    }

    for testName, test := range testCases {
        test := test
        t.Run(testName, func(t *testing.T) {
            t.Parallel()
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUseCase := mock.NewMockUserUseCase(ctrl)
            test.buildStub(mockUseCase, test.input)

            userHandler := handlers.NewUserHandler(mockUseCase)

            app := fiber.New()
            app.Post("/signin", userHandler.UserSignIn)

            jsonData, err := json.Marshal(test.input)
            assert.NoError(t, err)
            body := bytes.NewBuffer(jsonData)

            req := httptest.NewRequest("POST", "/signin", body)
            req.Header.Set("Content-Type", "application/json")
            resp, err := app.Test(req, -1)
            assert.NoError(t, err)
            test.checkResponse(t, resp)
        })
    }
}


func Test_UserSignOut(t *testing.T) {
	t.Run("Successful User SignOut", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mock.NewMockUserUseCase(ctrl)
		userHandler := handlers.NewUserHandler(mockUseCase)

		app := fiber.New()
		app.Post("/signout", userHandler.UserSignOut)

		req := httptest.NewRequest("POST", "/signout", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
