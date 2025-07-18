package model

import "errors"

type TaskRepository interface {
	CreateTask() TaskID
	AddURL(TaskID, string) error
	Status(TaskID) error
}

var (
	ErrTaskNotFound = errors.New("task not found")
)
