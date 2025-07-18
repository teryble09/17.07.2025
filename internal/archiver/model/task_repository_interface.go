package model

import "errors"

type TaskRepository interface {
	CreateTask() TaskID
	AddURL(TaskID, string) error
	Status(TaskID) ([]Url, error)
}

var (
	ErrTaskNotFound             = errors.New("task not found")
	ErrMaximumTaskNumberReached = errors.New("maximum task number reached")
)
