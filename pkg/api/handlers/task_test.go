package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"taskmanagementapi/pkg/api/handlers"
	"taskmanagementapi/pkg/usecase/mock"
	"taskmanagementapi/pkg/utils/models"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateTask(t *testing.T) {
    testCases := map[string]struct {
        input         models.CreateTask
        userID        string
        buildStub     func(useCaseMock *mock.MockTaskUseCase, task models.CreateTask, userID string)
        checkResponse func(t *testing.T, resp *http.Response)
    }{
        "Valid Task Creation": {
            input: models.CreateTask{
                Title:       "New Task",
                Description: "This is a new task.",
            },
            userID: "1",
            buildStub: func(useCaseMock *mock.MockTaskUseCase, task models.CreateTask, userID string) {
                useCaseMock.EXPECT().CreateTask(task, userID).Times(1).Return(nil)
            },
            checkResponse: func(t *testing.T, resp *http.Response) {
                assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
            },
        },
        "Invalid Task Input": {
            input: models.CreateTask{
                Title:       "",
                Description: "",
            },
            userID: "1",
            buildStub: func(useCaseMock *mock.MockTaskUseCase, task models.CreateTask, userID string) {
            },
            checkResponse: func(t *testing.T, resp *http.Response) {
                assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
            },
        },
        "Task Creation Failure": {
            input: models.CreateTask{
                Title:       "New Task",
                Description: "This is a new task.",
            },
            userID: "1",
            buildStub: func(useCaseMock *mock.MockTaskUseCase, task models.CreateTask, userID string) {
                useCaseMock.EXPECT().CreateTask(task, userID).Times(1).Return(errors.New("task creation failed"))
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

            mockUseCase := mock.NewMockTaskUseCase(ctrl)
            test.buildStub(mockUseCase, test.input, test.userID)

            taskHandler := handlers.NewTaskHandler(mockUseCase)

            app := fiber.New()
            app.Post("/task", func(c *fiber.Ctx) error {
                c.Locals("user_id", test.userID)
                return taskHandler.CreateTask(c)
            })
            jsonData, err := json.Marshal(test.input)
            assert.NoError(t, err)
            body := bytes.NewBuffer(jsonData)

            req := httptest.NewRequest("POST", "/task", body)
            req.Header.Set("Content-Type", "application/json")
            resp, err := app.Test(req, -1)
            assert.NoError(t, err)
            test.checkResponse(t, resp)
        })
    }
}


func Test_GetTasks(t *testing.T) {
	testCases := map[string]struct {
		userID        string
		buildStub     func(useCaseMock *mock.MockTaskUseCase, userID string)
		checkResponse func(t *testing.T, resp *http.Response)
	}{
		"Successfully Retrieve Tasks": {
			userID: "1",
			buildStub: func(useCaseMock *mock.MockTaskUseCase, userID string) {
				createdAt, _ := time.Parse(time.RFC3339, "2024-10-08T00:28:52+05:30")
				useCaseMock.EXPECT().GetTasks(userID).Times(1).Return([]models.TaskDetails{
					{
						ID:          "6705824f80a09eb0313f0e4",
						Title:       "Task 1",
						Description: "This is task 1",
						CreatedAt:   createdAt,
					},
					{
						ID:          "6705824f80a09eb0313f0e4",
						Title:       "Task 2",
						Description: "This is task 2",
						CreatedAt:   createdAt,
					},
				}, nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			},
		},
		"Tasks Retrieval Failure": {
			userID: "6705824f80a09eb0313f0e42",
			buildStub: func(useCaseMock *mock.MockTaskUseCase, userID string) {
				useCaseMock.EXPECT().GetTasks(userID).Times(1).Return(nil, errors.New("failed to get tasks"))
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

			mockUseCase := mock.NewMockTaskUseCase(ctrl)
			test.buildStub(mockUseCase, test.userID)

			taskHandler := handlers.NewTaskHandler(mockUseCase)

			app := fiber.New()
			app.Get("/tasks", func(c *fiber.Ctx) error {
				c.Locals("user_id", test.userID)
				return taskHandler.GetTasks(c)
			})

			req := httptest.NewRequest("GET", "/tasks", nil)
			req.Header.Set("Authorization", "Bearer some-token")
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			test.checkResponse(t, resp)
		})
	}
}
func Test_GetTask(t *testing.T) {
    testCases := map[string]struct {
        userID        string
        taskID        string
        buildStub     func(useCaseMock *mock.MockTaskUseCase, userID, taskID string)
        checkResponse func(t *testing.T, resp *http.Response)
    }{
        "Successfully Retrieve Task": {
            userID: "6705824f80a09eb0313f0e42",
            taskID: "6705824f80a09eb0313f0e4",
            buildStub: func(useCaseMock *mock.MockTaskUseCase, userID, taskID string) {
                createdAt, _ := time.Parse(time.RFC3339, "2024-10-08T00:28:52+05:30")
                useCaseMock.EXPECT().GetTask(userID, taskID).Times(1).Return(models.TaskDetails{
                    ID:          taskID,
                    Title:       "Task 1",
                    Description: "This is task 1",
                    CreatedAt:   createdAt,
                }, nil) 
            },
            checkResponse: func(t *testing.T, resp *http.Response) {
                assert.Equal(t, fiber.StatusOK, resp.StatusCode)
            },
        },
        "Task Retrieval Failure": {
            userID: "6705824f80a09eb0313f0e42",
            taskID: "6705824f80a09eb0313f0e4",
            buildStub: func(useCaseMock *mock.MockTaskUseCase, userID, taskID string) {
                useCaseMock.EXPECT().GetTask(userID, taskID).Times(1).Return(models.TaskDetails{}, errors.New("failed to get task"))
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

            mockUseCase := mock.NewMockTaskUseCase(ctrl)
            test.buildStub(mockUseCase, test.userID, test.taskID)

            taskHandler := handlers.NewTaskHandler(mockUseCase)

            app := fiber.New()
            app.Get("/tasks/:id", func(c *fiber.Ctx) error {
                c.Locals("user_id", test.userID)
                return taskHandler.GetTask(c)
            })

            req := httptest.NewRequest("GET", "/tasks/"+test.taskID, nil)
            req.Header.Set("Authorization", "Bearer some-token")
            resp, err := app.Test(req, -1)
            assert.NoError(t, err)
            test.checkResponse(t, resp)
        })
    }
}



func Test_UpdateTask(t *testing.T) {
	testCases := map[string]struct {
		userID        string
		taskID        string
		input         models.CreateTask
		buildStub     func(useCaseMock *mock.MockTaskUseCase, userID, taskID string, task models.CreateTask)
		checkResponse func(t *testing.T, resp *http.Response)
	}{
		"Successfully Update Task": {
			userID: "1",
			taskID: "1",
			input: models.CreateTask{
				Title:       "Updated Task",
				Description: "This is an updated task.",
			},
			buildStub: func(useCaseMock *mock.MockTaskUseCase, userID, taskID string, task models.CreateTask) {
				useCaseMock.EXPECT().UpdateTask(userID, taskID, task).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			},
		},
		"Invalid Task Input": {
			userID: "1",
			taskID: "1",
			input: models.CreateTask{
				Title:       "",
				Description: "",
			},
			buildStub: func(useCaseMock *mock.MockTaskUseCase, userID, taskID string, task models.CreateTask) {
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			},
		},
		"Task Update Failure": {
			userID: "1",
			taskID: "1",
			input: models.CreateTask{
				Title:       "Updated Task",
				Description: "This is an updated task.",
			},
			buildStub: func(useCaseMock *mock.MockTaskUseCase, userID, taskID string, task models.CreateTask) {
				useCaseMock.EXPECT().UpdateTask(userID, taskID, task).Times(1).Return(errors.New("update failed"))
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

			mockUseCase := mock.NewMockTaskUseCase(ctrl)
			test.buildStub(mockUseCase, test.userID, test.taskID, test.input)

			taskHandler := handlers.NewTaskHandler(mockUseCase)

			app := fiber.New()
			app.Put("/task/:id", func(c *fiber.Ctx) error {
				c.Locals("user_id", test.userID)
				return taskHandler.UpdateTask(c)
			})

			jsonData, err := json.Marshal(test.input)
			require.NoError(t, err)
			body := bytes.NewBuffer(jsonData)

			req := httptest.NewRequest("PUT", "/task/"+test.taskID, body)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			test.checkResponse(t, resp)
		})
	}
}

func Test_DeleteTask(t *testing.T) {
	testCases := map[string]struct {
		userID        string
		taskID        string
		buildStub     func(useCaseMock *mock.MockTaskUseCase, userID, taskID string)
		checkResponse func(t *testing.T, resp *http.Response)
	}{
		"Successfully Delete Task": {
			userID: "1",
			taskID: "1",
			buildStub: func(useCaseMock *mock.MockTaskUseCase, userID, taskID string) {
				useCaseMock.EXPECT().DeleteTask(userID, taskID).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			},
		},
		"Task Deletion Failure": {
			userID: "1",
			taskID: "1",
			buildStub: func(useCaseMock *mock.MockTaskUseCase, userID, taskID string) {
				useCaseMock.EXPECT().DeleteTask(userID, taskID).Times(1).Return(errors.New("task deletion failed"))
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

			mockUseCase := mock.NewMockTaskUseCase(ctrl)
			test.buildStub(mockUseCase, test.userID, test.taskID)

			taskHandler := handlers.NewTaskHandler(mockUseCase)

			app := fiber.New()
			app.Delete("/task/:id", func(c *fiber.Ctx) error {
				c.Locals("user_id", test.userID)
				return taskHandler.DeleteTask(c)
			})

			req := httptest.NewRequest("DELETE", "/task/"+test.taskID, nil)
			req.Header.Set("Authorization", "Bearer some-token")
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			test.checkResponse(t, resp)
		})
	}
}
