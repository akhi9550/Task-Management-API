package di

import (
	server "taskmanagementapi/pkg/api"
	"taskmanagementapi/pkg/api/handlers"
	"taskmanagementapi/pkg/config"
	"taskmanagementapi/pkg/db"
	"taskmanagementapi/pkg/repository"
	"taskmanagementapi/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*server.ServerHTTP, error) {
	database, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(database)
	taskRepository := repository.NewTaskRepository(database)

	UserUseCase := usecase.NewUserUseCase(userRepository)
	TaskUseCase := usecase.NewTaskUseCase(taskRepository)

	userHandler := handlers.NewUserHandler(UserUseCase)
	taskHandler := handlers.NewTaskHandler(TaskUseCase)

	serverHttp := server.NewServerHTTP(userHandler,
		taskHandler,
	)

	return serverHttp, nil
}
