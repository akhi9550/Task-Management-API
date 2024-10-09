run:  ##run code
	go run ./cmd/main.go

test: ##test
	go test ./... 

mock: ##make mock files using mockgen
	mockgen -source pkg\repository\interface\user.go -destination pkg\repository\mock\user_mock.go -package mock
	mockgen -source pkg\repository\interface\task.go -destination pkg\repository\mock\task_mock.go -package mock
	mockgen -source pkg\usecase\interface\user.go -destination pkg\usecase\mock\user_mock.go -package mock
	mockgen -source pkg\usecase\interface\task.go -destination pkg\usecase\mock\task_mock.go -package mock
	mockgen -source go.mongodb.org\mongo-driver\mongo -destination pkg\repository\mongomock\mongo_mock.go -package=mock