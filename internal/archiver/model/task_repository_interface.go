package model

type TaskRepository interface {
	CreateTask() TaskID
	AddURL(TaskID) error
	Status(TaskID)
}
