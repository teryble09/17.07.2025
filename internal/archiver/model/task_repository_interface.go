package model

type TaskRepository interface {
	CreateTask() string
	AddURL(string) error
	Status(string)
}
